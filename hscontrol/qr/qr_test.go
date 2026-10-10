package qr

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// PNG signature: first 8 bytes of a valid PNG file.
var pngSignature = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}

func TestGenerateRegistrationQR_Valid(t *testing.T) {
	authID := "hskey-abc123"
	serverURL := "https://headscale.example.com"
	expiresAt := time.Now().Add(24 * time.Hour)

	pngData, err := GenerateRegistrationQR(authID, serverURL, expiresAt)

	require.NoError(t, err)
	assert.NotEmpty(t, pngData)

	// Verify PNG signature
	assert.True(
		t,
		bytes.HasPrefix(pngData, pngSignature),
		"generated data should have valid PNG signature",
	)
}

func TestGenerateRegistrationQR_EmptyAuthID(t *testing.T) {
	serverURL := "https://headscale.example.com"
	expiresAt := time.Now().Add(24 * time.Hour)

	pngData, err := GenerateRegistrationQR("", serverURL, expiresAt)

	require.Error(t, err)
	assert.Nil(t, pngData)
	assert.Contains(t, err.Error(), "auth_id cannot be empty")
}

func TestGenerateRegistrationQR_EmptyServerURL(t *testing.T) {
	authID := "hskey-abc123"
	expiresAt := time.Now().Add(24 * time.Hour)

	pngData, err := GenerateRegistrationQR(authID, "", expiresAt)

	require.Error(t, err)
	assert.Nil(t, pngData)
	assert.Contains(t, err.Error(), "server_url cannot be empty")
}

func TestGenerateRegistrationQR_PastExpiry(t *testing.T) {
	authID := "hskey-abc123"
	serverURL := "https://headscale.example.com"
	expiresAt := time.Now().Add(-24 * time.Hour)

	pngData, err := GenerateRegistrationQR(authID, serverURL, expiresAt)

	require.Error(t, err)
	assert.Nil(t, pngData)
	assert.Contains(t, err.Error(), "expires_at must be in the future")
}

func TestGenerateRegistrationQRWithSize_CustomSize(t *testing.T) {
	authID := "hskey-abc123"
	serverURL := "https://headscale.example.com"
	expiresAt := time.Now().Add(24 * time.Hour)
	customSize := 512

	pngData, err := GenerateRegistrationQRWithSize(
		authID,
		serverURL,
		expiresAt,
		customSize,
	)

	require.NoError(t, err)
	assert.NotEmpty(t, pngData)

	// Verify PNG signature
	assert.True(
		t,
		bytes.HasPrefix(pngData, pngSignature),
		"generated data should have valid PNG signature",
	)

	// Custom size should generate a larger PNG
	defaultPng, err := GenerateRegistrationQR(authID, serverURL, expiresAt)
	require.NoError(t, err)
	assert.Greater(
		t,
		len(pngData),
		len(defaultPng),
		"larger size should generate more bytes",
	)
}

func TestGenerateRegistrationQRWithSize_SmallSize(t *testing.T) {
	authID := "hskey-abc123"
	serverURL := "https://headscale.example.com"
	expiresAt := time.Now().Add(24 * time.Hour)
	smallSize := 128

	pngData, err := GenerateRegistrationQRWithSize(
		authID,
		serverURL,
		expiresAt,
		smallSize,
	)

	require.NoError(t, err)
	assert.NotEmpty(t, pngData)

	// Verify PNG signature
	assert.True(
		t,
		bytes.HasPrefix(pngData, pngSignature),
		"generated data should have valid PNG signature",
	)
}

func TestGenerateRegistrationQR_LongPayload(t *testing.T) {
	// Test with a longer payload to ensure QR can handle it
	authID := "hskey-very-long-authentication-id-with-many-characters-123456789"
	serverURL := "https://very-long-headscale-server-url.example.com/with/path"
	expiresAt := time.Now().Add(24 * time.Hour)

	pngData, err := GenerateRegistrationQR(authID, serverURL, expiresAt)

	require.NoError(t, err)
	assert.NotEmpty(t, pngData)
	assert.True(
		t,
		bytes.HasPrefix(pngData, pngSignature),
		"generated data should have valid PNG signature",
	)
}

func TestGenerateRegistrationQR_RoundTrip(t *testing.T) {
	authID := "hskey-abc123"
	serverURL := "https://headscale.example.com"
	expiresAt := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)

	// Generate QR code
	pngData, err := GenerateRegistrationQR(authID, serverURL, expiresAt)
	require.NoError(t, err)

	// Verify it's a valid PNG
	assert.True(
		t,
		bytes.HasPrefix(pngData, pngSignature),
		"generated data should have valid PNG signature",
	)

	// Note: Full round-trip (decode QR, parse JSON) would require
	// an image decoder and QR decoder library, which is out of scope
	// for this unit test. The integration test or manual testing
	// should verify the QR code can be scanned and decoded correctly.
}

// BenchmarkGenerateRegistrationQR measures QR code generation performance.
// Target: < 100ms per operation.
func BenchmarkGenerateRegistrationQR(b *testing.B) {
	authID := "hskey-abc123"
	serverURL := "https://headscale.example.com"
	expiresAt := time.Now().Add(24 * time.Hour)

	b.ResetTimer()

	for range b.N {
		_, err := GenerateRegistrationQR(authID, serverURL, expiresAt)
		if err != nil {
			b.Fatalf("failed to generate QR: %v", err)
		}
	}
}

// BenchmarkGenerateRegistrationQRWithSize_Large benchmarks larger QR codes.
func BenchmarkGenerateRegistrationQRWithSize_Large(b *testing.B) {
	authID := "hskey-abc123"
	serverURL := "https://headscale.example.com"
	expiresAt := time.Now().Add(24 * time.Hour)
	size := 512

	b.ResetTimer()

	for range b.N {
		_, err := GenerateRegistrationQRWithSize(
			authID,
			serverURL,
			expiresAt,
			size,
		)
		if err != nil {
			b.Fatalf("failed to generate QR: %v", err)
		}
	}
}

func TestGenerateRegistrationQR_ZeroExpiry(t *testing.T) {
	authID := "hskey-abc123"
	serverURL := "https://headscale.example.com"
	expiresAt := time.Time{} // Zero time

	pngData, err := GenerateRegistrationQR(authID, serverURL, expiresAt)

	require.Error(t, err)
	assert.Nil(t, pngData)
	assert.Contains(t, err.Error(), "expires_at cannot be zero")
}

func TestGenerateRegistrationQR_ExpiryOnBoundary(t *testing.T) {
	// Test with expiry just barely in the future (1 millisecond)
	authID := "hskey-abc123"
	serverURL := "https://headscale.example.com"
	expiresAt := time.Now().Add(1 * time.Millisecond)

	pngData, err := GenerateRegistrationQR(authID, serverURL, expiresAt)

	// Should succeed if called before expiry
	require.NoError(t, err)
	assert.NotEmpty(t, pngData)
	assert.True(
		t,
		bytes.HasPrefix(pngData, pngSignature),
		"generated data should have valid PNG signature",
	)
}

func TestGenerateRegistrationQR_VeryShortExpiry(t *testing.T) {
	// Test that a QR generated with very short expiry works immediately
	// but would fail if regenerated after expiry
	authID := "hskey-abc123"
	serverURL := "https://headscale.example.com"
	expiresAt := time.Now().Add(10 * time.Millisecond)

	// First generation should succeed
	pngData1, err := GenerateRegistrationQR(authID, serverURL, expiresAt)
	require.NoError(t, err)
	assert.NotEmpty(t, pngData1)

	var (
		pngData2      []byte
		generationErr error
	)

	assert.EventuallyWithT(t, func(c *assert.CollectT) {
		pngData2, generationErr = GenerateRegistrationQR(authID, serverURL, expiresAt)
		assert.Error(c, generationErr)
	}, time.Second, time.Millisecond)

	assert.Nil(t, pngData2)
	assert.ErrorIs(t, generationErr, errExpiryNotInFuture)
}
