package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"

	"github.com/blackhiden/headscale/hscontrol/types"
	"gorm.io/gorm"
)

var (
	ErrWebhookNotFound      = errors.New("webhook not found")
	ErrWebhookNameExists    = errors.New("webhook with this name already exists")
	ErrWebhookInvalidEvents = errors.New("webhook events are invalid")
)

// CreateWebhook creates a new webhook in the database.
func (hsdb *HSDatabase) CreateWebhook(webhook *types.Webhook) error {
	// Validate events before saving
	err := validateWebhookEvents(webhook.Events)
	if err != nil {
		return err
	}

	// Check for duplicate name
	var existing types.Webhook

	err = hsdb.DB.Where("name = ?", webhook.Name).First(&existing).Error
	if err == nil {
		return ErrWebhookNameExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("checking for duplicate webhook name: %w", err)
	}

	err = hsdb.DB.Create(webhook).Error
	if err != nil {
		return fmt.Errorf("creating webhook: %w", err)
	}

	return nil
}

// GetWebhook retrieves a webhook by ID.
func (hsdb *HSDatabase) GetWebhook(id uint) (*types.Webhook, error) {
	var webhook types.Webhook

	err := hsdb.DB.First(&webhook, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrWebhookNotFound
		}

		return nil, fmt.Errorf("getting webhook: %w", err)
	}

	return &webhook, nil
}

// ListWebhooks returns all webhooks.
func (hsdb *HSDatabase) ListWebhooks() ([]*types.Webhook, error) {
	var webhooks []*types.Webhook

	err := hsdb.DB.Order("name ASC").Find(&webhooks).Error
	if err != nil {
		return nil, fmt.Errorf("listing webhooks: %w", err)
	}

	return webhooks, nil
}

// ListWebhooksForEvent returns all enabled webhooks that listen to a specific event.
func (hsdb *HSDatabase) ListWebhooksForEvent(eventType types.WebhookEventType) ([]*types.Webhook, error) {
	var webhooks []*types.Webhook

	// Get all enabled webhooks
	err := hsdb.DB.Where("enabled = ?", true).Find(&webhooks).Error
	if err != nil {
		return nil, fmt.Errorf("listing webhooks: %w", err)
	}

	// Filter by event type
	var filtered []*types.Webhook

	for _, wh := range webhooks {
		var events []types.WebhookEventType

		err := json.Unmarshal([]byte(wh.Events), &events)
		if err != nil {
			// Skip malformed webhooks, log elsewhere
			continue
		}

		if slices.Contains(events, eventType) {
			filtered = append(filtered, wh)
		}
	}

	return filtered, nil
}

// UpdateWebhook updates an existing webhook.
func (hsdb *HSDatabase) UpdateWebhook(webhook *types.Webhook) error {
	// Validate events before saving
	err := validateWebhookEvents(webhook.Events)
	if err != nil {
		return err
	}

	// Check if exists
	var existing types.Webhook

	err = hsdb.DB.First(&existing, webhook.ID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrWebhookNotFound
		}

		return fmt.Errorf("checking webhook existence: %w", err)
	}

	// Check for duplicate name (excluding self)
	var duplicate types.Webhook

	err = hsdb.DB.Where("name = ? AND id != ?", webhook.Name, webhook.ID).First(&duplicate).Error
	if err == nil {
		return ErrWebhookNameExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("checking for duplicate webhook name: %w", err)
	}

	err = hsdb.DB.Save(webhook).Error
	if err != nil {
		return fmt.Errorf("updating webhook: %w", err)
	}

	return nil
}

// DeleteWebhook deletes a webhook by ID.
func (hsdb *HSDatabase) DeleteWebhook(id uint) error {
	result := hsdb.DB.Delete(&types.Webhook{}, id)
	if result.Error != nil {
		return fmt.Errorf("deleting webhook: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrWebhookNotFound
	}

	return nil
}

// validateWebhookEvents validates that the events JSON contains valid event types.
func validateWebhookEvents(eventsJSON string) error {
	var events []types.WebhookEventType

	err := json.Unmarshal([]byte(eventsJSON), &events)
	if err != nil {
		return fmt.Errorf("invalid events JSON: %w", err)
	}

	if len(events) == 0 {
		return fmt.Errorf("%w: at least one event type required", ErrWebhookInvalidEvents)
	}

	for _, e := range events {
		if !e.IsValid() {
			return fmt.Errorf("%w: invalid event type '%s'", ErrWebhookInvalidEvents, e)
		}
	}

	return nil
}
