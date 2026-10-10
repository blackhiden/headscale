package apiv1

import (
	"context"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/arsydoni4326-alt/headscale/hscontrol/updatecheck"
)

// updateCheckInput mirrors the query parameters of the update-check endpoint.
// check is a string so the exact "true" (case-insensitive) comparison of the
// underlying handler is preserved.
type updateCheckInput struct {
	Check string `doc:"When \"true\", additionally fetch the remote latest commit and compare it with the running binary." query:"check"`
}

type updateCheckOutput struct {
	Body updatecheck.UpdateCheckResult
}

func init() {
	registrations = append(registrations, registerUpdateCheck)
}

// registerUpdateCheck exposes the fork's update-check endpoint as a v1 API
// operation so it is discoverable in the OpenAPI document. It is intentionally
// unauthenticated (no Security requirement), matching the endpoint's no-auth
// behavior. The server's disable_check_updates config option is wired through
// so operators can opt out of the remote comparison.
func registerUpdateCheck(api huma.API, b Backend) {
	huma.Register(api, huma.Operation{
		OperationID: "update-check",
		Method:      http.MethodGet,
		Path:        "/api/v1/update-check",
		Summary:     "Check for updates",
		Description: "Returns the running binary's version information. When ?check=true is set, additionally fetches the latest commit from the fork's main branch and reports whether an update is available.",
		Tags:        []string{"Update"},
	}, func(ctx context.Context, input *updateCheckInput) (*updateCheckOutput, error) {
		return &updateCheckOutput{
			Body: updatecheck.BuildResponse(strings.EqualFold(input.Check, "true"), b.Cfg.DisableUpdateCheck),
		}, nil
	})
}
