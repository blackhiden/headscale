package qr

import (
	"fmt"
	"time"

	"github.com/skip2/go-qrcode"
)

const (
	// DefaultQRSize is the default pixel size for generated QR codes.
	DefaultQRSize = 256
	// DefaultRecoveryLevel is the default error correction level.
	DefaultRecoveryLevel = qrcode.Medium
)

// GenerateRegistrationQR generates a QR code PNG image for node registration.
// The QR code encodes a JSON payload with the provided authentication ID,
// server URL, and expiry time.
//
// Returns PNG image bytes or an error if generation fails.
func GenerateRegistrationQR(
	authID, serverURL string,
	expiresAt time.Time,
) ([]byte, error) {
	return GenerateRegistrationQRWithSize(
		authID,
		serverURL,
		expiresAt,
		DefaultQRSize,
	)
}

// GenerateRegistrationQRWithSize generates a QR code PNG image with a custom
// pixel size. The QR code encodes a JSON payload with the provided
// authentication ID, server URL, and expiry time.
//
// Returns PNG image bytes or an error if generation fails.
func GenerateRegistrationQRWithSize(
	authID, serverURL string,
	expiresAt time.Time,
	size int,
) ([]byte, error) {
	// Create and validate the payload
	payload, err := NewRegistrationPayload(authID, serverURL, expiresAt)
	if err != nil {
		return nil, fmt.Errorf("invalid payload: %w", err)
	}

	// Serialize to JSON
	jsonData, err := payload.ToJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize payload: %w", err)
	}

	// Generate QR code
	qr, err := qrcode.New(string(jsonData), DefaultRecoveryLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to create QR code: %w", err)
	}

	// Encode to PNG with the specified size
	pngData, err := qr.PNG(size)
	if err != nil {
		return nil, fmt.Errorf("failed to encode QR code as PNG: %w", err)
	}

	return pngData, nil
}
