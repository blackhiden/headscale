package hscontrol

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/blackhiden/headscale/hscontrol/qr"
	"github.com/blackhiden/headscale/hscontrol/types"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// RegistrationQRPayloadHandler returns the QR registration payload as JSON
// for a given auth_id. This is an open endpoint (no authentication required)
// that allows CLI clients to fetch the payload and generate QR codes locally.
//
// GET /api/v1/registration/qr/{auth_id}
//
// Returns:
//   - 200 OK: JSON payload with type, version, auth_id, server_url, expires_at
//   - 400 Bad Request: Invalid auth_id format
//   - 404 Not Found: Pending registration not found
//   - 410 Gone: Registration has expired
func (h *Headscale) RegistrationQRPayloadHandler(
	writer http.ResponseWriter,
	req *http.Request,
) {
	authIDStr := chi.URLParam(req, "auth_id")
	if authIDStr == "" {
		http.Error(writer, "missing auth_id parameter", http.StatusBadRequest)
		return
	}

	authID, err := types.AuthIDFromString(authIDStr)
	if err != nil {
		log.Debug().
			Str("auth_id", authIDStr).
			Err(err).
			Msg("invalid auth_id format")
		http.Error(writer, "invalid auth_id format", http.StatusBadRequest)
		return
	}

	// Lookup pending registration
	authReq, ok := h.state.GetAuthCacheEntry(authID)
	if !ok {
		log.Debug().
			Str("auth_id", authID.String()).
			Msg("pending registration not found")
		http.Error(writer, "pending registration not found", http.StatusNotFound)
		return
	}

	// Verify this is a registration request (not SSH check or other auth type)
	if !authReq.IsRegistration() {
		log.Debug().
			Str("auth_id", authID.String()).
			Msg("auth request is not a registration request")
		http.Error(writer, "auth request is not a registration request", http.StatusNotFound)
		return
	}

	// Check if expired
	if authReq.ExpiresAt().Before(time.Now()) {
		log.Debug().
			Str("auth_id", authID.String()).
			Time("expires_at", authReq.ExpiresAt()).
			Msg("registration has expired")
		http.Error(writer, "registration has expired", http.StatusGone)
		return
	}

	// Build the payload using the existing qr package
	payload, err := qr.NewRegistrationPayload(
		authID.String(),
		h.cfg.ServerURL,
		authReq.ExpiresAt(),
	)
	if err != nil {
		log.Error().
			Err(err).
			Str("auth_id", authID.String()).
			Msg("failed to create QR payload")
		http.Error(writer, "internal server error", http.StatusInternalServerError)
		return
	}

	// Serialize to JSON
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	
	if err := json.NewEncoder(writer).Encode(payload); err != nil {
		log.Error().
			Err(err).
			Str("auth_id", authID.String()).
			Msg("failed to encode QR payload to JSON")
	}
}
