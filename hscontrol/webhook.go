package hscontrol

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/arsydoni4326-alt/headscale/hscontrol/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog/log"
)

var (
	webhookDispatchTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: prometheusNamespace,
		Name:      "webhook_dispatch_total",
		Help:      "Total webhook dispatches by event type and status",
	}, []string{"event_type", "status"})

	webhookDispatchDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: prometheusNamespace,
		Name:      "webhook_dispatch_duration_seconds",
		Help:      "Webhook dispatch duration in seconds",
		Buckets:   prometheus.DefBuckets,
	}, []string{"event_type"})
)

// WebhookDispatcher handles sending webhook notifications for events.
type WebhookDispatcher struct {
	app *Headscale
}

// NewWebhookDispatcher creates a new webhook dispatcher.
func NewWebhookDispatcher(app *Headscale) *WebhookDispatcher {
	return &WebhookDispatcher{
		app: app,
	}
}

// DispatchNodeUp sends notifications when a node comes online.
func (wd *WebhookDispatcher) DispatchNodeUp(node types.NodeView) {
	data := types.NodeEventData{
		NodeID:   node.ID().Uint64(),
		NodeName: node.GivenName(),
		UserName: node.User().Name(),
	}

	wd.dispatch(types.WebhookEventNodeUp, data)
}

// DispatchNodeDown sends notifications when a node goes offline.
func (wd *WebhookDispatcher) DispatchNodeDown(node types.NodeView) {
	data := types.NodeEventData{
		NodeID:   node.ID().Uint64(),
		NodeName: node.GivenName(),
		UserName: node.User().Name(),
	}

	wd.dispatch(types.WebhookEventNodeDown, data)
}

// DispatchHealthCheckFail sends notifications when a health check fails.
func (wd *WebhookDispatcher) DispatchHealthCheckFail(endpoint, errorMsg string) {
	data := types.HealthCheckEventData{
		Endpoint: endpoint,
		Error:    errorMsg,
	}

	wd.dispatch(types.WebhookEventHealthCheckFail, data)
}

// DispatchAlertFiring sends notifications for external alerts.
func (wd *WebhookDispatcher) DispatchAlertFiring(alertName, severity string, labels map[string]string, details map[string]interface{}) {
	data := types.AlertEventData{
		AlertName: alertName,
		Severity:  severity,
		Labels:    labels,
		Details:   details,
	}

	wd.dispatch(types.WebhookEventAlertFiring, data)
}

// dispatch is the internal method that sends webhook requests.
func (wd *WebhookDispatcher) dispatch(eventType types.WebhookEventType, data interface{}) {
	timer := prometheus.NewTimer(webhookDispatchDuration.WithLabelValues(string(eventType)))
	defer timer.ObserveDuration()

	// Get webhooks for this event type
	webhooks, err := wd.app.state.ListWebhooksForEvent(eventType)
	if err != nil {
		log.Error().
			Err(err).
			Str("event_type", string(eventType)).
			Msg("Failed to list webhooks for event")
		webhookDispatchTotal.WithLabelValues(string(eventType), "list_error").Inc()
		return
	}

	if len(webhooks) == 0 {
		// No webhooks configured for this event type, nothing to do
		return
	}

	// Create payload
	payload := types.WebhookPayload{
		EventType: eventType,
		Timestamp: time.Now().UTC(),
		Data:      data,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		log.Error().
			Err(err).
			Str("event_type", string(eventType)).
			Msg("Failed to marshal webhook payload")
		webhookDispatchTotal.WithLabelValues(string(eventType), "marshal_error").Inc()
		return
	}

	// Send to each webhook asynchronously
	for _, webhook := range webhooks {
		go wd.sendWebhook(webhook, payloadJSON, eventType)
	}
}

// sendWebhook sends a single webhook request.
func (wd *WebhookDispatcher) sendWebhook(webhook *types.Webhook, payloadJSON []byte, eventType types.WebhookEventType) {
	timeout := time.Duration(webhook.TimeoutSeconds) * time.Second
	if timeout == 0 {
		timeout = 10 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create request
	req, err := http.NewRequestWithContext(ctx, "POST", webhook.URL, bytes.NewReader(payloadJSON))
	if err != nil {
		log.Error().
			Err(err).
			Uint("webhook_id", webhook.ID).
			Str("webhook_name", webhook.Name).
			Str("event_type", string(eventType)).
			Msg("Failed to create webhook request")
		webhookDispatchTotal.WithLabelValues(string(eventType), "request_error").Inc()
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Headscale-Webhook/1.0")

	// Add custom headers
	if webhook.Headers != "" {
		var headers map[string]string
		if err := json.Unmarshal([]byte(webhook.Headers), &headers); err != nil {
			log.Error().
				Err(err).
				Uint("webhook_id", webhook.ID).
				Str("webhook_name", webhook.Name).
				Msg("Failed to unmarshal webhook headers")
		} else {
			for k, v := range headers {
				req.Header.Set(k, v)
			}
		}
	}

	// Add HMAC signature if secret is set
	if webhook.Secret != "" {
		signature := computeHMAC(payloadJSON, webhook.Secret)
		req.Header.Set("X-Headscale-Signature", signature)
	}

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().
			Err(err).
			Uint("webhook_id", webhook.ID).
			Str("webhook_name", webhook.Name).
			Str("webhook_url", webhook.URL).
			Str("event_type", string(eventType)).
			Msg("Failed to send webhook")
		webhookDispatchTotal.WithLabelValues(string(eventType), "send_error").Inc()
		return
	}
	defer resp.Body.Close()

	// Read response body for logging (limit to 1KB)
	bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Info().
			Uint("webhook_id", webhook.ID).
			Str("webhook_name", webhook.Name).
			Str("event_type", string(eventType)).
			Int("status_code", resp.StatusCode).
			Msg("Webhook sent successfully")
		webhookDispatchTotal.WithLabelValues(string(eventType), "success").Inc()
	} else {
		log.Warn().
			Uint("webhook_id", webhook.ID).
			Str("webhook_name", webhook.Name).
			Str("webhook_url", webhook.URL).
			Str("event_type", string(eventType)).
			Int("status_code", resp.StatusCode).
			Str("response_body", string(bodyBytes)).
			Msg("Webhook returned non-2xx status")
		webhookDispatchTotal.WithLabelValues(string(eventType), "http_error").Inc()
	}
}

// computeHMAC computes an HMAC-SHA256 signature of the payload.
func computeHMAC(payload []byte, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(payload)
	return hex.EncodeToString(h.Sum(nil))
}
