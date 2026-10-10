package apiv1

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/danielgtaylor/huma/v2"
	"github.com/arsydoni4326-alt/headscale/hscontrol/types"
)

func init() {
	registrations = append(registrations, registerWebhooks)
}

// Webhook is the v1 Webhook message.
type Webhook struct {
	ID             uint     `json:"id"`
	Name           string   `json:"name"`
	URL            string   `json:"url"`
	Events         []string `json:"events"`
	Enabled        bool     `json:"enabled"`
	TimeoutSeconds int      `json:"timeoutSeconds"`
	CreatedAt      string   `json:"createdAt"`
	UpdatedAt      string   `json:"updatedAt"`
}

// CreateWebhookRequestBody is the request body for creating a webhook.
type CreateWebhookRequestBody struct {
	Name           string            `json:"name" minLength:"1" maxLength:"255" doc:"Human-readable webhook identifier"`
	URL            string            `json:"url" format:"uri" doc:"Destination endpoint for webhook POST requests"`
	Events         []string          `json:"events" minItems:"1" doc:"Event types that trigger this webhook"`
	Headers        map[string]string `json:"headers,omitempty" doc:"Custom HTTP headers to include in requests"`
	Secret         string            `json:"secret,omitempty" doc:"Secret for HMAC-SHA256 payload signing"`
	Enabled        *bool             `json:"enabled,omitempty" doc:"Whether this webhook is active (default: true)"`
	TimeoutSeconds *int              `json:"timeoutSeconds,omitempty" doc:"HTTP timeout in seconds (default: 10)"`
}

// UpdateWebhookRequestBody is the request body for updating a webhook.
type UpdateWebhookRequestBody struct {
	Name           *string           `json:"name,omitempty" minLength:"1" maxLength:"255"`
	URL            *string           `json:"url,omitempty" format:"uri"`
	Events         []string          `json:"events,omitempty" minItems:"1"`
	Headers        map[string]string `json:"headers,omitempty"`
	Secret         *string           `json:"secret,omitempty"`
	Enabled        *bool             `json:"enabled,omitempty"`
	TimeoutSeconds *int              `json:"timeoutSeconds,omitempty"`
}

type (
	createWebhookInput struct {
		Body CreateWebhookRequestBody
	}
	createWebhookOutput struct {
		Body struct {
			Webhook Webhook `json:"webhook"`
		}
	}
)

type (
	getWebhookInput struct {
		ID string `path:"id" format:"uint" doc:"Webhook ID"`
	}
	getWebhookOutput struct {
		Body struct {
			Webhook Webhook `json:"webhook"`
		}
	}
)

type (
	listWebhooksOutput struct {
		Body struct {
			Webhooks []Webhook `json:"webhooks" nullable:"false"`
		}
	}
)

type (
	updateWebhookInput struct {
		ID   string `path:"id" format:"uint" doc:"Webhook ID"`
		Body UpdateWebhookRequestBody
	}
	updateWebhookOutput struct {
		Body struct {
			Webhook Webhook `json:"webhook"`
		}
	}
)

type (
	deleteWebhookInput struct {
		ID string `path:"id" format:"uint" doc:"Webhook ID"`
	}
	deleteWebhookOutput struct {
		Body struct{}
	}
)

func registerWebhooks(api huma.API, b Backend) {
	huma.Register(api, huma.Operation{
		OperationID: "createWebhook",
		Method:      http.MethodPost,
		Path:        "/api/v1/webhooks",
		Summary:     "Create webhook",
		Description: "Creates a new webhook configuration for monitoring and alerting integrations",
		Tags:        []string{"Webhooks"},
		Security:    bearerAuth,
	}, func(ctx context.Context, in *createWebhookInput) (*createWebhookOutput, error) {
		// Validate events
		for _, eventStr := range in.Body.Events {
			eventType := types.WebhookEventType(eventStr)
			if !eventType.IsValid() {
				return nil, huma.Error400BadRequest("invalid event type: " + eventStr)
			}
		}

		// Marshal events to JSON
		eventsJSON, err := json.Marshal(in.Body.Events)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid events", err)
		}

		// Marshal headers to JSON
		headersJSON := ""
		if in.Body.Headers != nil {
			headersBytes, err := json.Marshal(in.Body.Headers)
			if err != nil {
				return nil, huma.Error400BadRequest("invalid headers", err)
			}
			headersJSON = string(headersBytes)
		}

		// Set defaults
		enabled := true
		if in.Body.Enabled != nil {
			enabled = *in.Body.Enabled
		}

		timeoutSeconds := 10
		if in.Body.TimeoutSeconds != nil {
			timeoutSeconds = *in.Body.TimeoutSeconds
		}

		webhook := &types.Webhook{
			Name:           in.Body.Name,
			URL:            in.Body.URL,
			Events:         string(eventsJSON),
			Headers:        headersJSON,
			Secret:         in.Body.Secret,
			Enabled:        enabled,
			TimeoutSeconds: timeoutSeconds,
		}

		err = b.State.CreateWebhook(webhook)
		if err != nil {
			return nil, mapError("creating webhook", err)
		}

		out := &createWebhookOutput{}
		out.Body.Webhook = webhookFromState(webhook)

		return out, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "getWebhook",
		Method:      http.MethodGet,
		Path:        "/api/v1/webhooks/{id}",
		Summary:     "Get webhook",
		Description: "Retrieves a webhook configuration by ID",
		Tags:        []string{"Webhooks"},
		Security:    bearerAuth,
	}, func(ctx context.Context, in *getWebhookInput) (*getWebhookOutput, error) {
		id, err := parseWebhookID(in.ID)
		if err != nil {
			return nil, err
		}

		webhook, err := b.State.GetWebhook(id)
		if err != nil {
			return nil, mapError("getting webhook", err)
		}

		out := &getWebhookOutput{}
		out.Body.Webhook = webhookFromState(webhook)

		return out, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "listWebhooks",
		Method:      http.MethodGet,
		Path:        "/api/v1/webhooks",
		Summary:     "List webhooks",
		Description: "Lists all webhook configurations",
		Tags:        []string{"Webhooks"},
		Security:    bearerAuth,
	}, func(ctx context.Context, _ *struct{}) (*listWebhooksOutput, error) {
		webhooks, err := b.State.ListWebhooks()
		if err != nil {
			return nil, huma.Error500InternalServerError("listing webhooks", err)
		}

		out := &listWebhooksOutput{}
		out.Body.Webhooks = make([]Webhook, len(webhooks))
		for i, wh := range webhooks {
			out.Body.Webhooks[i] = webhookFromState(wh)
		}

		return out, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "updateWebhook",
		Method:      http.MethodPut,
		Path:        "/api/v1/webhooks/{id}",
		Summary:     "Update webhook",
		Description: "Updates an existing webhook configuration",
		Tags:        []string{"Webhooks"},
		Security:    bearerAuth,
	}, func(ctx context.Context, in *updateWebhookInput) (*updateWebhookOutput, error) {
		id, err := parseWebhookID(in.ID)
		if err != nil {
			return nil, err
		}

		webhook, err := b.State.GetWebhook(id)
		if err != nil {
			return nil, mapError("getting webhook", err)
		}

		// Update fields if provided
		if in.Body.Name != nil {
			webhook.Name = *in.Body.Name
		}
		if in.Body.URL != nil {
			webhook.URL = *in.Body.URL
		}
		if in.Body.Events != nil {
			// Validate events
			for _, eventStr := range in.Body.Events {
				eventType := types.WebhookEventType(eventStr)
				if !eventType.IsValid() {
					return nil, huma.Error400BadRequest("invalid event type: " + eventStr)
				}
			}

			eventsJSON, err := json.Marshal(in.Body.Events)
			if err != nil {
				return nil, huma.Error400BadRequest("invalid events", err)
			}
			webhook.Events = string(eventsJSON)
		}
		if in.Body.Headers != nil {
			headersJSON, err := json.Marshal(in.Body.Headers)
			if err != nil {
				return nil, huma.Error400BadRequest("invalid headers", err)
			}
			webhook.Headers = string(headersJSON)
		}
		if in.Body.Secret != nil {
			webhook.Secret = *in.Body.Secret
		}
		if in.Body.Enabled != nil {
			webhook.Enabled = *in.Body.Enabled
		}
		if in.Body.TimeoutSeconds != nil {
			webhook.TimeoutSeconds = *in.Body.TimeoutSeconds
		}

		err = b.State.UpdateWebhook(webhook)
		if err != nil {
			return nil, mapError("updating webhook", err)
		}

		out := &updateWebhookOutput{}
		out.Body.Webhook = webhookFromState(webhook)

		return out, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "deleteWebhook",
		Method:      http.MethodDelete,
		Path:        "/api/v1/webhooks/{id}",
		Summary:     "Delete webhook",
		Description: "Deletes a webhook configuration",
		Tags:        []string{"Webhooks"},
		Security:    bearerAuth,
	}, func(ctx context.Context, in *deleteWebhookInput) (*deleteWebhookOutput, error) {
		id, err := parseWebhookID(in.ID)
		if err != nil {
			return nil, err
		}

		err = b.State.DeleteWebhook(id)
		if err != nil {
			return nil, mapError("deleting webhook", err)
		}

		return &deleteWebhookOutput{}, nil
	})
}

// parseWebhookID decodes the webhook ID from a string.
func parseWebhookID(s string) (uint, error) {
	id, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, huma.Error400BadRequest("invalid webhook id", err)
	}

	return uint(id), nil
}

// webhookFromState converts a domain webhook into the v1 response shape.
func webhookFromState(wh *types.Webhook) Webhook {
	// Unmarshal events
	var events []string
	_ = json.Unmarshal([]byte(wh.Events), &events)

	return Webhook{
		ID:             wh.ID,
		Name:           wh.Name,
		URL:            wh.URL,
		Events:         events,
		Enabled:        wh.Enabled,
		TimeoutSeconds: wh.TimeoutSeconds,
		CreatedAt:      wh.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      wh.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
