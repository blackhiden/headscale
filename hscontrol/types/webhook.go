package types

import (
	"time"

	"gorm.io/gorm"
)

// WebhookEventType represents the type of event that triggers a webhook.
type WebhookEventType string

const (
	// WebhookEventNodeUp is fired when a node comes online.
	WebhookEventNodeUp WebhookEventType = "node_up"
	// WebhookEventNodeDown is fired when a node goes offline.
	WebhookEventNodeDown WebhookEventType = "node_down"
	// WebhookEventHealthCheckFail is fired when a health check fails.
	WebhookEventHealthCheckFail WebhookEventType = "health_check_fail"
	// WebhookEventAlertFiring is fired when an external alert is received.
	WebhookEventAlertFiring WebhookEventType = "alert_firing"
)

// AllWebhookEventTypes returns all valid webhook event types.
func AllWebhookEventTypes() []WebhookEventType {
	return []WebhookEventType{
		WebhookEventNodeUp,
		WebhookEventNodeDown,
		WebhookEventHealthCheckFail,
		WebhookEventAlertFiring,
	}
}

// IsValid returns true if the event type is valid.
func (e WebhookEventType) IsValid() bool {
	for _, valid := range AllWebhookEventTypes() {
		if e == valid {
			return true
		}
	}
	return false
}

// Webhook represents a webhook configuration in the database.
type Webhook struct {
	gorm.Model

	// Name is a human-readable identifier for the webhook.
	Name string `gorm:"uniqueIndex;not null"`

	// URL is the destination endpoint for webhook POST requests.
	URL string `gorm:"not null"`

	// Events is a JSON array of event types that trigger this webhook.
	// Stored as JSON to avoid a separate join table.
	Events string `gorm:"type:text;not null"`

	// Headers is a JSON object of custom HTTP headers to include in requests.
	// Use for authentication tokens or custom metadata.
	Headers string `gorm:"type:text"`

	// Secret is used to sign webhook payloads with HMAC-SHA256.
	// The signature is sent in the X-Headscale-Signature header.
	// Optional; if empty, payloads are not signed.
	Secret string

	// Enabled controls whether this webhook is active.
	Enabled bool `gorm:"not null;default:true"`

	// TimeoutSeconds is the HTTP timeout for webhook requests.
	// Defaults to 10 seconds if not set.
	TimeoutSeconds int `gorm:"not null;default:10"`
}

// WebhookPayload is the structure sent to webhook endpoints.
type WebhookPayload struct {
	EventType WebhookEventType `json:"event_type"`
	Timestamp time.Time        `json:"timestamp"`
	Data      interface{}      `json:"data"`
}

// NodeEventData contains node-specific event data.
type NodeEventData struct {
	NodeID   uint64 `json:"node_id"`
	NodeName string `json:"node_name"`
	UserName string `json:"user_name"`
}

// HealthCheckEventData contains health check failure data.
type HealthCheckEventData struct {
	Endpoint string `json:"endpoint"`
	Error    string `json:"error"`
}

// AlertEventData contains generic alert data.
type AlertEventData struct {
	AlertName string                 `json:"alert_name"`
	Severity  string                 `json:"severity"`
	Labels    map[string]string      `json:"labels,omitempty"`
	Details   map[string]interface{} `json:"details,omitempty"`
}
