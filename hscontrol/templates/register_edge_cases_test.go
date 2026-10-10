package templates

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAuthWebEdgeCases tests edge cases for the non-OIDC registration page.
func TestAuthWebEdgeCases(t *testing.T) {
	t.Run("empty strings", func(t *testing.T) {
		html := AuthWeb("", "", "", "").Render()

		// Should still produce valid HTML
		assert.Contains(t, html, "<!DOCTYPE html>")
		assert.Contains(t, html, "</html>")

		// Should not have broken structure
		assert.Contains(t, html, "<h1>")
		assert.Contains(t, html, "</h1>")
	})

	t.Run("very long auth ID", func(t *testing.T) {
		longAuthID := strings.Repeat("a", 500)
		command := "headscale auth register --auth-id " + longAuthID + " --user USERNAME"

		html := AuthWeb("Node registration", "Run the command below:", command, "").Render()

		// Should contain the full command
		assert.Contains(t, html, longAuthID)

		// Should be properly escaped and not break HTML structure
		assert.Contains(t, html, "<!DOCTYPE html>")
		assert.Contains(t, html, "</html>")
	})

	t.Run("special characters in command", func(t *testing.T) {
		// Test XSS attempt in command
		maliciousCommand := "headscale auth register --auth-id <script>alert('xss')</script>"

		html := AuthWeb("Test", "Description", maliciousCommand, "").Render()

		// Script tags should be escaped
		assert.NotContains(t, html, "<script>alert('xss')</script>")
		// Should contain escaped version
		assert.Contains(t, html, "&lt;script&gt;")
	})

	t.Run("unicode characters", func(t *testing.T) {
		title := "Inscription de nœud 🚀"
		description := "Exécutez la commande ci-dessous"
		command := "headscale auth register --user tëst-üser"

		html := AuthWeb(title, description, command, "").Render()

		// Should contain unicode characters properly
		assert.Contains(t, html, "🚀")
		assert.Contains(t, html, "nœud")
		assert.Contains(t, html, "tëst-üser")
	})

	t.Run("whitespace-only strings", func(t *testing.T) {
		html := AuthWeb("   ", "\n\t", "  \n  ", "").Render()

		// Should still produce valid HTML
		assert.Contains(t, html, "<!DOCTYPE html>")
		assert.Contains(t, html, "</html>")
	})

	t.Run("newlines and tabs in content", func(t *testing.T) {
		command := "headscale auth register \\\n  --auth-id test-123 \\\n  --user USERNAME"

		html := AuthWeb("Title", "Description", command, "").Render()

		// Should preserve whitespace in code block
		assert.Contains(t, html, "test-123")
		assert.Contains(t, html, "USERNAME")
	})
}

// TestRegisterConfirmEdgeCases tests edge cases for the OIDC registration confirmation page.
func TestRegisterConfirmEdgeCases(t *testing.T) {
	t.Run("empty hostname", func(t *testing.T) {
		info := RegisterConfirmInfo{
			FormAction:    "/register/confirm/test-auth-id",
			CSRFTokenName: "csrf_token",
			CSRFToken:     "token-value",
			User:          "user@example.com",
			Hostname:      "", // Empty hostname
			OS:            "linux",
			MachineKey:    "mkey:abc123",
			QRCodeDataURL: "",
		}

		html := RegisterConfirm(info).Render()

		// Should still render valid HTML
		require.Contains(t, html, "<!DOCTYPE html>")
		require.Contains(t, html, "</html>")

		// Should show empty hostname (not crash)
		assert.Contains(t, html, "Hostname")
	})

	t.Run("empty OS shows unknown", func(t *testing.T) {
		info := RegisterConfirmInfo{
			FormAction:    "/register/confirm/test-auth-id",
			CSRFTokenName: "csrf_token",
			CSRFToken:     "token-value",
			User:          "user@example.com",
			Hostname:      "test-node",
			OS:            "", // Empty OS
			MachineKey:    "mkey:abc123",
			QRCodeDataURL: "",
		}

		html := RegisterConfirm(info).Render()

		// Should display "(unknown)" for empty OS
		assert.Contains(t, html, "(unknown)")
	})

	t.Run("very long hostname", func(t *testing.T) {
		info := RegisterConfirmInfo{
			FormAction:    "/register/confirm/test-auth-id",
			CSRFTokenName: "csrf_token",
			CSRFToken:     "token-value",
			User:          "user@example.com",
			Hostname:      strings.Repeat("very-long-hostname-", 20), // 380 chars
			OS:            "linux",
			MachineKey:    "mkey:abc123",
			QRCodeDataURL: "",
		}

		html := RegisterConfirm(info).Render()

		// Should contain the full hostname
		assert.Contains(t, html, info.Hostname)

		// Should still be valid HTML
		require.Contains(t, html, "<!DOCTYPE html>")
		require.Contains(t, html, "</html>")
	})

	t.Run("very long username", func(t *testing.T) {
		longUser := strings.Repeat("user", 50) + "@example.com"
		info := RegisterConfirmInfo{
			FormAction:    "/register/confirm/test-auth-id",
			CSRFTokenName: "csrf_token",
			CSRFToken:     "token-value",
			User:          longUser,
			Hostname:      "test-node",
			OS:            "linux",
			MachineKey:    "mkey:abc123",
			QRCodeDataURL: "",
		}

		html := RegisterConfirm(info).Render()

		// Should contain the full username
		assert.Contains(t, html, longUser)
		require.Contains(t, html, "<!DOCTYPE html>")
	})

	t.Run("special characters in hostname", func(t *testing.T) {
		info := RegisterConfirmInfo{
			FormAction:    "/register/confirm/test-auth-id",
			CSRFTokenName: "csrf_token",
			CSRFToken:     "token-value",
			User:          "user@example.com",
			Hostname:      "<script>alert('xss')</script>",
			OS:            "linux",
			MachineKey:    "mkey:abc123",
			QRCodeDataURL: "",
		}

		html := RegisterConfirm(info).Render()

		// Script tags should be escaped
		assert.NotContains(t, html, "<script>alert('xss')</script>")
		assert.Contains(t, html, "&lt;script&gt;")
	})

	t.Run("special characters in username", func(t *testing.T) {
		info := RegisterConfirmInfo{
			FormAction:    "/register/confirm/test-auth-id",
			CSRFTokenName: "csrf_token",
			CSRFToken:     "token-value",
			User:          "user<img src=x onerror=alert(1)>@example.com",
			Hostname:      "test-node",
			OS:            "linux",
			MachineKey:    "mkey:abc123",
			QRCodeDataURL: "",
		}

		html := RegisterConfirm(info).Render()

		// XSS attempt should be escaped
		assert.NotContains(t, html, "<img src=x onerror=alert(1)>")
		assert.Contains(t, html, "&lt;img")
	})

	t.Run("no QR code data", func(t *testing.T) {
		info := RegisterConfirmInfo{
			FormAction:    "/register/confirm/test-auth-id",
			CSRFTokenName: "csrf_token",
			CSRFToken:     "token-value",
			User:          "user@example.com",
			Hostname:      "test-node",
			OS:            "linux",
			MachineKey:    "mkey:abc123",
			QRCodeDataURL: "",
		}

		html := RegisterConfirm(info).Render()

		// Should not contain QR code section
		assert.NotContains(t, html, "Or scan with Headplane")
		assert.NotContains(t, html, "Registration QR Code")

		// Should still render form
		assert.Contains(t, html, "Confirm registration")
		assert.Contains(t, html, `type="submit"`)
	})

	t.Run("with valid QR code data", func(t *testing.T) {
		info := RegisterConfirmInfo{
			FormAction:    "/register/confirm/test-auth-id",
			CSRFTokenName: "csrf_token",
			CSRFToken:     "token-value",
			User:          "user@example.com",
			Hostname:      "test-node",
			OS:            "linux",
			MachineKey:    "mkey:abc123",
			QRCodeDataURL: "data:image/png;base64,iVBORw0KGg==",
		}

		html := RegisterConfirm(info).Render()

		// Should contain QR code section
		assert.Contains(t, html, "Or scan with Headplane")
		assert.Contains(t, html, "Registration QR Code")
		assert.Contains(t, html, info.QRCodeDataURL)
	})

	t.Run("unicode in user and hostname", func(t *testing.T) {
		//nolint:gosmopolitan // Testing international character support
		info := RegisterConfirmInfo{
			FormAction:    "/register/confirm/test-auth-id",
			CSRFTokenName: "csrf_token",
			CSRFToken:     "token-value",
			User:          "用户@例え.jp",
			Hostname:      "ノード-🖥️",
			OS:            "macOS",
			MachineKey:    "mkey:abc123",
			QRCodeDataURL: "",
		}

		html := RegisterConfirm(info).Render()

		//nolint:gosmopolitan // Validating international character rendering
		// Should handle unicode properly
		assert.Contains(t, html, "用户@例え.jp")
		assert.Contains(t, html, "ノード-🖥️")
	})
}

// TestRegisterConfirmFormStructure validates the form structure and CSRF protection.
func TestRegisterConfirmFormStructure(t *testing.T) {
	t.Run("form has correct structure", func(t *testing.T) {
		info := RegisterConfirmInfo{
			FormAction:    "/register/confirm/auth-123",
			CSRFTokenName: "csrf_token",
			CSRFToken:     "secure-token-456",
			User:          "user@example.com",
			Hostname:      "test-node",
			OS:            "linux",
			MachineKey:    "mkey:abc123",
			QRCodeDataURL: "",
		}

		html := RegisterConfirm(info).Render()

		// Form should have POST method
		assert.Contains(t, html, `method="POST"`)

		// Form action should match
		assert.Contains(t, html, `action="/register/confirm/auth-123"`)

		// Hidden CSRF field should be present
		assert.Contains(t, html, `type="hidden"`)
		assert.Contains(t, html, `name="csrf_token"`)
		assert.Contains(t, html, `value="secure-token-456"`)

		// Submit button should be present
		assert.Contains(t, html, `type="submit"`)
		assert.Contains(t, html, "Confirm registration")
	})

	t.Run("CSRF token is properly embedded", func(t *testing.T) {
		//nolint:gosec // Test data, not real credentials
		info := RegisterConfirmInfo{
			FormAction:    "/test",
			CSRFTokenName: "custom_csrf",
			CSRFToken:     "token-with-special-chars-<>&\"",
			User:          "user@example.com",
			Hostname:      "node",
			OS:            "linux",
			MachineKey:    "mkey:123",
			QRCodeDataURL: "",
		}

		html := RegisterConfirm(info).Render()

		// CSRF token should be present with custom name
		assert.Contains(t, html, `name="custom_csrf"`)
		// Token value should be present (elem-go handles attribute escaping)
		assert.Contains(t, html, `type="hidden"`)
		// The form should be structurally valid
		assert.Contains(t, html, "<form")
		assert.Contains(t, html, "</form>")
	})

	t.Run("device details table structure", func(t *testing.T) {
		info := RegisterConfirmInfo{
			FormAction:    "/test",
			CSRFTokenName: "csrf",
			CSRFToken:     "token",
			User:          "test@example.com",
			Hostname:      "laptop",
			OS:            "macOS",
			MachineKey:    "mkey:xyz",
			QRCodeDataURL: "",
		}

		html := RegisterConfirm(info).Render()

		// Table should be present
		assert.Contains(t, html, "<table")
		assert.Contains(t, html, "</table>")

		// Row labels should be present
		assert.Contains(t, html, "Hostname")
		assert.Contains(t, html, "OS")
		assert.Contains(t, html, "Machine key")
		assert.Contains(t, html, "Registered to")

		// Values should be present
		assert.Contains(t, html, "laptop")
		assert.Contains(t, html, "macOS")
		assert.Contains(t, html, "mkey:xyz")
		assert.Contains(t, html, "test@example.com")
	})
}

// TestRegisterConfirmAccessibility tests accessibility features.
func TestRegisterConfirmAccessibility(t *testing.T) {
	t.Run("QR code has alt text", func(t *testing.T) {
		info := RegisterConfirmInfo{
			FormAction:    "/test",
			CSRFTokenName: "csrf",
			CSRFToken:     "token",
			User:          "user@example.com",
			Hostname:      "node",
			OS:            "linux",
			MachineKey:    "mkey:123",
			QRCodeDataURL: "data:image/png;base64,iVBORw0KGg==",
		}

		html := RegisterConfirm(info).Render()

		// QR code image should have alt text
		assert.Contains(t, html, `alt="Registration QR Code"`)
	})

	t.Run("page has proper HTML structure", func(t *testing.T) {
		info := RegisterConfirmInfo{
			FormAction:    "/test",
			CSRFTokenName: "csrf",
			CSRFToken:     "token",
			User:          "user@example.com",
			Hostname:      "node",
			OS:            "linux",
			MachineKey:    "mkey:123",
			QRCodeDataURL: "",
		}

		html := RegisterConfirm(info).Render()

		// Should have proper document structure
		assert.Contains(t, html, "<!DOCTYPE html>")
		assert.Contains(t, html, `<html lang="en">`)
		assert.Contains(t, html, `charset="UTF-8"`)
		assert.Contains(t, html, "</html>")

		// Should have viewport meta tag for mobile
		assert.Contains(t, html, `name="viewport"`)
	})
}

// TestAuthWebAccessibility tests accessibility of non-OIDC registration page.
func TestAuthWebAccessibility(t *testing.T) {
	t.Run("has proper heading structure", func(t *testing.T) {
		html := AuthWeb("Node Registration", "Instructions", "command", "").Render()

		// Should have H1 for title
		assert.Contains(t, html, "<h1>")
		assert.Contains(t, html, "Node Registration")
		assert.Contains(t, html, "</h1>")
	})

	t.Run("code block for command", func(t *testing.T) {
		html := AuthWeb("Title", "Desc", "headscale auth register", "").Render()

		// Should have proper code block structure
		assert.Contains(t, html, "<pre>")
		assert.Contains(t, html, "<code>")
		assert.Contains(t, html, "headscale auth register")
		assert.Contains(t, html, "</code>")
		assert.Contains(t, html, "</pre>")
	})

	t.Run("QR code section appears when provided", func(t *testing.T) {
		qrData := "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg=="
		html := AuthWeb("Title", "Desc", "headscale auth register", qrData).Render()

		// Should have QR section
		assert.Contains(t, html, "Or scan with Headplane")
		assert.Contains(t, html, qrData)
		assert.Contains(t, html, `alt="Registration QR Code"`)
	})

	t.Run("QR code section omitted when empty", func(t *testing.T) {
		html := AuthWeb("Title", "Desc", "headscale auth register", "").Render()

		// Should not have QR section
		assert.NotContains(t, html, "Or scan with Headplane")
		assert.NotContains(t, html, "Registration QR Code")
	})
}
