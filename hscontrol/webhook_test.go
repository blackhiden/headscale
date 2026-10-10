package hscontrol

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/arsydoni4326-alt/headscale/hscontrol/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComputeHMAC(t *testing.T) {
	payload := []byte(`{"event_type":"node_up","timestamp":"2026-09-30T00:00:00Z"}`)
	secret := "test-secret"

	signature := computeHMAC(payload, secret)

	// Verify signature is a valid hex string
	assert.Len(t, signature, 64) // SHA256 produces 64 hex chars
	assert.Regexp(t, "^[0-9a-f]{64}$", signature)

	// Verify signature is deterministic
	signature2 := computeHMAC(payload, secret)
	assert.Equal(t, signature, signature2)

	// Verify different payloads produce different signatures
	differentPayload := []byte(`{"event_type":"node_down","timestamp":"2026-09-30T00:00:00Z"}`)
	differentSignature := computeHMAC(differentPayload, secret)
	assert.NotEqual(t, signature, differentSignature)
}

func TestWebhookDispatcher_SendWebhook(t *testing.T) {
	// Create a test HTTP server to receive webhook
	receivedPayload := ""
	receivedHeaders := http.Header{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Headscale-Webhook/1.0", r.Header.Get("User-Agent"))

		// Store headers for verification
		receivedHeaders = r.Header.Clone()

		// Read body
		buf := make([]byte, 1024)
		n, _ := r.Body.Read(buf)
		receivedPayload = string(buf[:n])

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create webhook with secret
	eventsJSON, _ := json.Marshal([]string{"node_up"})
	webhook := &types.Webhook{
		ID:             1,
		Name:           "test-webhook",
		URL:            server.URL,
		Events:         string(eventsJSON),
		Secret:         "test-secret",
		Enabled:        true,
		TimeoutSeconds: 5,
	}

	// Create a minimal app for testing
	cfg := &types.Config{}
	app, err := NewHeadscale(cfg)
	require.NoError(t, err)

	dispatcher := NewWebhookDispatcher(app)

	// Prepare payload
	payload := types.WebhookPayload{
		EventType: types.WebhookEventNodeUp,
		Timestamp: time.Now().UTC(),
		Data: types.NodeEventData{
			NodeID:   123,
			NodeName: "test-node",
			UserName: "test-user",
		},
	}

	payloadJSON, err := json.Marshal(payload)
	require.NoError(t, err)

	// Send webhook
	dispatcher.sendWebhook(webhook, payloadJSON, types.WebhookEventNodeUp)

	// Give it a moment to complete
	time.Sleep(100 * time.Millisecond)

	// Verify payload was received
	assert.NotEmpty(t, receivedPayload)

	// Verify signature was sent
	signature := receivedHeaders.Get("X-Headscale-Signature")
	assert.NotEmpty(t, signature)

	// Verify signature is correct
	expectedSignature := computeHMAC(payloadJSON, "test-secret")
	assert.Equal(t, expectedSignature, signature)

	// Verify payload content
	var receivedWebhookPayload types.WebhookPayload
	err = json.Unmarshal([]byte(receivedPayload), &receivedWebhookPayload)
	require.NoError(t, err)
	assert.Equal(t, types.WebhookEventNodeUp, receivedWebhookPayload.EventType)
}

func TestWebhookEventType_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		event types.WebhookEventType
		valid bool
	}{
		{"node_up is valid", types.WebhookEventNodeUp, true},
		{"node_down is valid", types.WebhookEventNodeDown, true},
		{"health_check_fail is valid", types.WebhookEventHealthCheckFail, true},
		{"alert_firing is valid", types.WebhookEventAlertFiring, true},
		{"invalid event", types.WebhookEventType("invalid"), false},
		{"empty event", types.WebhookEventType(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.valid, tt.event.IsValid())
		})
	}
}

func TestAllWebhookEventTypes(t *testing.T) {
	events := types.AllWebhookEventTypes()

	assert.Len(t, events, 4)
	assert.Contains(t, events, types.WebhookEventNodeUp)
	assert.Contains(t, events, types.WebhookEventNodeDown)
	assert.Contains(t, events, types.WebhookEventHealthCheckFail)
	assert.Contains(t, events, types.WebhookEventAlertFiring)
}
