package templates

import (
	"cmp"

	"github.com/chasefleming/elem-go"
	"github.com/chasefleming/elem-go/attrs"
	"github.com/chasefleming/elem-go/styles"
)

// RegisterConfirmInfo carries the human-readable information shown on
// the registration confirmation interstitial that an OIDC-authenticated
// user must explicitly accept before a pending node is registered to
// their identity. The fields here intentionally include enough device
// detail (hostname, OS, machine-key fingerprint) for the user to
// recognise whether the device they are about to claim is in fact
// theirs.
type RegisterConfirmInfo struct {
	// FormAction is the absolute or relative URL the confirm form
	// POSTs to. Typically /register/confirm/{auth_id}.
	FormAction string

	// CSRFTokenName is the name of the hidden form field carrying the
	// CSRF token. The corresponding cookie shares this name.
	CSRFTokenName string

	// CSRFToken is the per-session token that must match the value of
	// the cookie set by the OIDC callback before the POST is honoured.
	CSRFToken string

	// User is the OIDC-authenticated identity the device will be
	// registered to if the user confirms.
	User string

	// Hostname is the hostname the registering tailscaled instance
	// reported in its [tailcfg.RegisterRequest].
	Hostname string

	// OS is the operating system the registering tailscaled reported.
	// May be the empty string when the client did not send [tailcfg.Hostinfo].
	OS string

	// MachineKey is the short fingerprint of the registering machine
	// key. The full key is intentionally not shown.
	MachineKey string

	// QRCodeDataURL is a base64-encoded data URL of the QR code PNG
	// image for Headplane mobile scanning. Empty string if QR generation
	// failed; the template omits the QR section when empty.
	QRCodeDataURL string
}

// RegisterConfirm renders an interstitial page that asks the
// OIDC-authenticated user to explicitly confirm that they want to
// register the named device under their account. Without this
// confirmation step a single GET to /register/{auth_id} could
// silently complete a phishing-style registration when the victim's
// IdP allows silent SSO.
func RegisterConfirm(info RegisterConfirmInfo) *elem.Element {
	deviceList := deviceTable(
		[]deviceRow{
			{"Hostname", elem.Text(info.Hostname)},
			{"OS", elem.Text(cmp.Or(info.OS, "(unknown)"))},
			{"Machine key", Code(elem.Text(info.MachineKey))},
			{"Registered to", elem.Text(info.User)},
		},
	)

	// Build content nodes: device table, optional QR section, then form
	content := []elem.Node{
		H2(elem.Text("Confirm node registration")),
		P(elem.Text(
			"A device is asking to be added to your tailnet. " +
				"Please review the details below and confirm that this device is yours.",
		)),
		deviceList,
	}

	// Add QR code section if available
	if info.QRCodeDataURL != "" {
		content = append(content, qrSection(info.QRCodeDataURL))
	}

	// Add confirmation form
	form := elem.Form(
		attrs.Props{
			attrs.Method: "POST",
			attrs.Action: info.FormAction,
		},
		elem.Input(attrs.Props{
			attrs.Type:  "hidden",
			attrs.Name:  info.CSRFTokenName,
			attrs.Value: info.CSRFToken,
		}),
		elem.Button(
			attrs.Props{attrs.Type: "submit"},
			elem.Text("Confirm registration"),
		),
	)
	content = append(content, form)

	// Add footer text
	content = append(content, P(elem.Text(
		"If you do not recognise this device, close this window. "+
			"The registration request will expire automatically.",
	)))

	return page("Headscale - Confirm node registration", content...)
}

type deviceRow struct {
	label string
	value elem.Node
}

func deviceTable(rows []deviceRow) *elem.Element {
	tableRows := make([]elem.Node, 0, len(rows))
	for _, row := range rows {
		tableRows = append(tableRows, elem.Tr(
			nil,
			elem.Td(attrs.Props{
				attrs.Style: styles.Props{
					styles.Padding:      "0.5rem 1rem 0.5rem 0",
					styles.FontWeight:   "600",
					styles.WhiteSpace:   "nowrap",
					styles.Color:        "var(--md-default-fg-color--light)",
					styles.BorderBottom: cssBorderHS,
				}.ToInline(),
			}, elem.Text(row.label)),
			elem.Td(attrs.Props{
				attrs.Style: styles.Props{
					styles.Padding:      "0.5rem 0",
					styles.BorderBottom: cssBorderHS,
				}.ToInline(),
			}, row.value),
		))
	}

	return elem.Table(attrs.Props{
		attrs.Style: styles.Props{
			styles.Width:          "100%",
			styles.BorderCollapse: "collapse",
			styles.MarginTop:      "1em",
			styles.MarginBottom:   "1.5em",
		}.ToInline(),
	}, tableRows...)
}

// qrSection creates a section displaying a QR code for Headplane mobile scanning.
// The QR code is embedded as a base64 data URL and styled to be responsive.
func qrSection(dataURL string) *elem.Element {
	return elem.Div(
		attrs.Props{
			attrs.Style: styles.Props{
				styles.MarginTop:    "2rem",
				styles.MarginBottom: "2rem",
				styles.Padding:      spaceL,
				styles.Background:   "var(--hs-bg)",
				styles.Border:       cssBorderHS,
				styles.BorderRadius: "0.5rem",
				styles.TextAlign:    cssCenter,
			}.ToInline(),
		},
		H3(elem.Text("Or scan with Headplane")),
		elem.P(
			attrs.Props{
				attrs.Style: styles.Props{
					styles.MarginBottom: spaceM,
					styles.Color:        "var(--md-default-fg-color--light)",
				}.ToInline(),
			},
			elem.Text("Navigate to: "),
			elem.Strong(nil, elem.Text("Machines → Scan QR")),
		),
		elem.Img(attrs.Props{
			attrs.Src: dataURL,
			attrs.Alt: "Registration QR Code",
			attrs.Style: styles.Props{
				styles.MaxWidth:     "256px",
				styles.Width:        "100%",
				styles.Height:       "auto",
				styles.Display:      "block",
				styles.MarginLeft:   "auto",
				styles.MarginRight:  "auto",
				styles.BorderRadius: "0.375rem",
			}.ToInline(),
		}),
	)
}
