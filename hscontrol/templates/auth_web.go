package templates

import (
	"github.com/chasefleming/elem-go"
)

// AuthWeb renders a page that instructs an administrator to run a CLI command
// to complete an authentication or registration flow.
// It is used by both the registration and auth-approve web handlers.
func AuthWeb(title, description, command, qrCodeDataURL string) *elem.Element {
	content := []elem.Node{
		H1(elem.Text(title)),
		P(elem.Text(description)),
		codeBlockText(command),
	}
	if qrCodeDataURL != "" {
		content = append(content, qrSection(qrCodeDataURL))
	}

	return page(title+" - Headscale", content...)
}
