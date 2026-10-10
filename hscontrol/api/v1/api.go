// Package apiv1 is the code-first Huma implementation of the Headscale v1 API.
// Handlers are a thin adapter over hscontrol/state; Huma emits the OpenAPI 3.1
// spec from the Go definitions (see Spec), and that spec drives the client.
//
// It depends only on the domain layer (hscontrol/state, hscontrol/types) via
// Backend, never on the hscontrol server package, so a future hscontrol/api/v2
// can sit beside it without either importing the other.
package apiv1

import (
	"context"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/arsydoni4326-alt/headscale/hscontrol/scope"
	"github.com/arsydoni4326-alt/headscale/hscontrol/state"
	"github.com/arsydoni4326-alt/headscale/hscontrol/types"
	"github.com/arsydoni4326-alt/headscale/hscontrol/types/change"
)

// Backend is the dependency surface the v1 API needs from the control plane:
// the state layer, the change-notification sink that distributes updates to
// connected nodes, and the config (only Policy.Mode and Policy.Path are read).
type Backend struct {
	State  *state.State
	Change func(...change.Change)
	Cfg    *types.Config
}

// NewAPI builds the v1 Huma API on the given chi router and registers every
// operation. Auth is enforced by a Huma middleware driven by each operation's
// declared bearer security (see authMiddleware); locally-trusted requests
// bypass it via WithLocalTrust.
func NewAPI(router chi.Router, backend Backend) huma.API {
	config := huma.DefaultConfig("Headscale API", "v1")
	config.Info.Description = "Headscale control server API."

	// Version the OpenAPI/docs routes under /api/v1 so a future v2 owns its own.
	// These register as plain mux routes, not operations, so they never appear
	// in the emitted spec or client.
	config.OpenAPIPath = "/api/v1/openapi"
	config.DocsPath = "/api/v1/docs"

	// The v1 API does not emit "$schema".
	config.SchemasPath = ""

	// Drop the default schema-link create hook: it injects a "$schema" property
	// and Link header into every response, which the v1 contract omits.
	config.CreateHooks = nil

	config.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"bearer": {
			Type:   "http",
			Scheme: "bearer",
		},
	}

	api := humachi.New(router, config)

	// Must run before register: Huma snapshots the middleware chain at operation
	// registration, so a middleware added afterwards would silently never run.
	api.UseMiddleware(authMiddleware(api, backend))

	register(api, backend)

	return api
}

// bearerAuth is the security requirement applied to every operation: all
// /api/v1 routes require an API key.
var bearerAuth = []map[string][]string{{"bearer": {}}}

// registrations is populated by each resource file's init(), so adding a
// resource group means adding a file rather than editing a shared point. Huma
// sorts the emitted spec, so init order does not affect output.
var registrations []func(huma.API, Backend)

// register wires up every operation contributed by the resource files.
func register(api huma.API, b Backend) {
	for _, fn := range registrations {
		fn(api, b)
	}
}

// Spec emits the OpenAPI 3.1 document. The zero Backend is safe because
// handlers are registered but never invoked during emission.
func Spec() ([]byte, error) {
	api := NewAPI(chi.NewMux(), Backend{})
	return api.OpenAPI().YAML()
}

// Spec30 emits the document downgraded to OpenAPI 3.0.3, needed because the
// client generator (oapi-codegen v2) cannot yet read the 3.1 spec.
func Spec30() ([]byte, error) {
	api := NewAPI(chi.NewMux(), Backend{})
	return api.OpenAPI().DowngradeYAML()
}

// Handler builds the v1 API on a fresh mux and returns both. Callers mount the
// mux and may use mux.Match to detect which paths this API serves.
func Handler(backend Backend) (*chi.Mux, huma.API) {
	mux := chi.NewMux()
	api := NewAPI(mux, backend)

	return mux, api
}

// localTrustKey marks a request as arriving over a locally-trusted transport;
// the auth middleware skips authentication for such requests.
type localTrustKey struct{}

// principalKey is the context key under which the auth middleware records the
// caller's principal: principalAdminKey is set for an all-access admin API key
// (or a locally-trusted request) and principalScopesKey carries the granted
// scopes of an OAuth access token. Handlers use these to apply capability-level
// RBAC beyond authentication, e.g. [registerMachines] requires the
// devices:core scope.
type principalKey int

const (
	principalAdminKey principalKey = iota
	principalScopesKey
)

// WithLocalTrust wraps a handler so its requests bypass API-key authentication.
// The unix socket uses this — access to the socket is the trust boundary — as
// do in-process tests that exercise the mux directly. Such requests are treated
// as all-access (admin), since the transport itself is privileged.
func WithLocalTrust(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		next.ServeHTTP(w, req.WithContext(
			context.WithValue(req.Context(), localTrustKey{}, struct{}{}),
		))
	})
}

// authMiddleware enforces the bearer credential for any operation that declares
// security and records the caller's principal on the context so handlers can
// apply capability-level RBAC. An admin API key is all-access (recorded as
// principalAdminKey); an OAuth access token is scope-limited (its granted
// scopes recorded under principalScopesKey) and rejected when it lacks the
// operation's declared scope (see [requireScope]). Locally-trusted requests and
// operations without declared security pass through. b.State is nil only during
// spec emission, where no request is served, so it is never dereferenced there.
func authMiddleware(api huma.API, b Backend) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		if ctx.Context().Value(localTrustKey{}) != nil {
			next(huma.WithValue(ctx, principalAdminKey, true))

			return
		}

		if len(ctx.Operation().Security) == 0 {
			next(ctx)

			return
		}

		token, ok := strings.CutPrefix(ctx.Header("Authorization"), "Bearer ")
		if !ok || token == "" {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "Unauthorized")

			return
		}

		// An OAuth access token is scope-limited; an admin API key is
		// all-access. They are told apart by prefix so a scoped token can never
		// be mistaken for an all-access key.
		if strings.HasPrefix(token, types.AccessTokenPrefix) {
			at, err := b.State.AuthenticateAccessToken(token)
			if err != nil {
				_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "Unauthorized")

				return
			}

			if want, ok := requiredScope(ctx.Operation()); ok && !scope.Grants(scope.Parse(at.Scopes), want) {
				_ = huma.WriteErr(api, ctx, http.StatusForbidden,
					"credential is missing the required scope "+string(want))

				return
			}

			next(huma.WithValue(ctx, principalScopesKey, at.Scopes))

			return
		}

		valid, err := b.State.ValidateAPIKey(token)
		if err != nil || !valid {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "Unauthorized")

			return
		}

		next(huma.WithValue(ctx, principalAdminKey, true))
	}
}

// isAdmin reports whether the request's principal holds all-access admin
// rights: an admin API key or a locally-trusted request. A scope-limited OAuth
// token is not admin.
func isAdmin(ctx context.Context) bool {
	admin, _ := ctx.Value(principalAdminKey).(bool)

	return admin
}

// principalScopes returns the scopes granted to the request's OAuth access
// token, and whether the request authenticated with one. ok is false for an
// admin API key, which is all-access and not scope-checked.
func principalScopes(ctx context.Context) ([]string, bool) {
	scopes, ok := ctx.Value(principalScopesKey).([]string)

	return scopes, ok
}

// scopeMetaKey keys the per-operation required scope in huma.Operation.Metadata.
const scopeMetaKey = "headscale.scope"

// requireScope records op's required scope, both in its Metadata (where the
// auth middleware reads it back) and in the generated OpenAPI document: an
// x-required-scope extension for machine consumers and a Description line so
// the rendered docs state what each operation needs. An admin API key is
// all-access and is never scope-checked.
func requireScope(op huma.Operation, s scope.Scope) huma.Operation {
	if op.Metadata == nil {
		op.Metadata = map[string]any{}
	}

	op.Metadata[scopeMetaKey] = s

	if op.Extensions == nil {
		op.Extensions = map[string]any{}
	}

	op.Extensions["x-required-scope"] = string(s)

	note := "Requires the `" + string(s) + "` OAuth scope (an admin API key is all-access)."
	if op.Description == "" {
		op.Description = note
	} else {
		op.Description += "\n\n" + note
	}

	return op
}

// requiredScope returns the scope an operation declared via requireScope, if any.
func requiredScope(op *huma.Operation) (scope.Scope, bool) {
	if op == nil || op.Metadata == nil {
		return "", false
	}

	s, ok := op.Metadata[scopeMetaKey].(scope.Scope)

	return s, ok
}
