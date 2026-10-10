// Package updatecheck provides an indestructible, self-contained update check
// endpoint for the headscale fork at github.com/arsydoni4326-alt/headscale.
//
// DO NOT REMOVE: This is a fork-specific feature that must survive upstream
// merges. All update-check logic lives in this single package so it can be
// restored as a unit if removed by an upstream merge.
//
// Endpoint: GET /api/v1/update-check
//
// Without ?check=true: returns current version info only (fast, no external call).
// With    ?check=true: fetches the remote latest release tag from the GitHub API,
//
//	compares it with the running binary's version when the local
//	build is a tagged release; falls back to commit-hash comparison
//	for dev builds. Results are cached for 15 minutes.
package updatecheck

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/arsydoni4326-alt/headscale/hscontrol/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog/log"
)

const (
	// DefaultRemoteRepoOwner is the default GitHub owner of the fork.
	DefaultRemoteRepoOwner = "arsydoni4326-alt"

	// DefaultRemoteRepoName is the default GitHub repository name of the fork.
	DefaultRemoteRepoName = "headscale"

	// RemoteRepo is the default fork repository URL for update checking.
	RemoteRepo = "https://github.com/arsydoni4326-alt/headscale.git"

	// RemoteRepoOwner is the default GitHub owner of the fork.
	//
	// Deprecated: kept for backward compatibility; use repoOwner() instead.
	RemoteRepoOwner = "arsydoni4326-alt"

	// RemoteRepoName is the default GitHub repository name of the fork.
	//
	// Deprecated: kept for backward compatibility; use repoName() instead.
	RemoteRepoName = "headscale"

	// remoteRepoEnv is the environment variable that overrides the remote
	// repository (format: owner/repo).
	remoteRepoEnv = "HEADSCALE_UPDATE_CHECK_REPO"

	// httpTimeout is the timeout for the GitHub API request.
	httpTimeout = 10 * time.Second

	// cacheTTL is how long a successful remote check result is cached.
	cacheTTL = 15 * time.Minute

	// unknownVersion is the fallback value for missing build metadata.
	unknownVersion = "unknown"
)

// Sentinel errors returned by the update check package.
var (
	errGitHubStatus   = errors.New("github api returned unexpected status")
	errGitHubEmptySHA = errors.New("github api returned empty sha")
	errGitHubEmptyTag = errors.New("github api returned empty tag_name")
	errNotAVersion    = errors.New("not a version string")
	errVersionParts   = errors.New("version does not have 3 dot-separated parts")
	errMajorNotNumber = errors.New("major version is not a number")
	errMinorNotNumber = errors.New("minor version is not a number")
	errPatchNotNumber = errors.New("patch version is not a number")
)

// CurrentVersionResponse is the version information for the running binary.
type CurrentVersionResponse struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildTime string `json:"buildTime"`
	Dirty     bool   `json:"dirty"`
}

// RemoteVersionResponse is the version information from the remote repository.
type RemoteVersionResponse struct {
	Commit  string `json:"commit"`
	Version string `json:"version,omitempty"`
	URL     string `json:"url"`
}

// UpdateCheckResult is the full response from the update check endpoint.
type UpdateCheckResult struct {
	Current         CurrentVersionResponse `json:"current"`
	UpdateAvailable *bool                  `json:"updateAvailable,omitempty"`
	Remote          *RemoteVersionResponse `json:"remote,omitempty"`
	Error           string                 `json:"error,omitempty"`
}

// semver represents a parsed semantic version (major.minor.patch-pre-release).
type semver struct {
	major, minor, patch int
	preRelease          string
}

// Prometheus metrics.
var (
	updateCheckRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "headscale",
		Name:      "updatecheck_requests_total",
		Help:      "Total number of update-check requests",
	}, []string{"check"})

	updateCheckRemoteFailures = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "headscale",
		Name:      "updatecheck_remote_failures_total",
		Help:      "Total number of failed remote update checks",
	}, []string{"reason"})

	updateCheckCacheHits = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "headscale",
		Name:      "updatecheck_cache_hits_total",
		Help:      "Total number of cache hits for remote update checks",
	})

	updateCheckCacheMisses = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "headscale",
		Name:      "updatecheck_cache_misses_total",
		Help:      "Total number of cache misses for remote update checks",
	})
)

// cacheEntry holds a cached remote check result.
type cacheEntry struct {
	value     string
	fetchedAt time.Time
}

var (
	cacheMu    sync.Mutex
	cache      = make(map[string]cacheEntry) // keyed by cache key ("commit" or "release")
	repoConfig sync.Once
	repoOwnerVal,
	repoNameVal string
)

// repoOwner returns the configured repo owner, defaulting to DefaultRemoteRepoOwner.
// Reads HEADSCALE_UPDATE_CHECK_REPO on first call.
func repoOwner() string {
	repoConfig.Do(initRepoConfig)
	return repoOwnerVal
}

// repoName returns the configured repo name, defaulting to DefaultRemoteRepoName.
func repoName() string {
	repoConfig.Do(initRepoConfig)
	return repoNameVal
}

func initRepoConfig() {
	repoOwnerVal = DefaultRemoteRepoOwner
	repoNameVal = DefaultRemoteRepoName

	repo := os.Getenv(remoteRepoEnv)
	if repo == "" {
		return
	}

	owner, name, ok := strings.Cut(repo, "/")
	if !ok || owner == "" || name == "" {
		log.Warn().
			Str("repo", repo).
			Str("env", remoteRepoEnv).
			Msg("update-check: invalid HEADSCALE_UPDATE_CHECK_REPO, expected owner/repo; using default")

		return
	}

	repoOwnerVal = owner
	repoNameVal = name
}

// repoURL returns the GitHub repository URL.
func repoURL() string {
	return fmt.Sprintf("https://github.com/%s/%s", repoOwner(), repoName())
}

// apiBaseURL is the base URL for the GitHub API. It is a variable so tests can
// point the fetch functions at a local httptest server.
var apiBaseURL = "https://api.github.com"

// commitAPIURL returns the GitHub API URL for the latest commit on main.
func commitAPIURL() string {
	return fmt.Sprintf("%s/repos/%s/%s/commits/main", apiBaseURL, repoOwner(), repoName())
}

// releaseAPIURL returns the GitHub API URL for the latest release.
func releaseAPIURL() string {
	return fmt.Sprintf("%s/repos/%s/%s/releases/latest", apiBaseURL, repoOwner(), repoName())
}

// BuildResponse builds the update check response.
//
// When check is true, it fetches the remote latest release tag (for release
// builds) or commit hash (for dev builds) from the fork's main branch and
// compares it with the running binary's embedded version info. Results are
// cached for cacheTTL to stay within GitHub's unauthenticated rate limits.
//
// When disabled is true (the upstream disable_check_updates config option),
// the remote check is skipped entirely and only the running binary's version
// info is returned, regardless of check. The v1 API operation wires this from
// the server config so operators can opt out of the remote comparison.
//
// It is the single source of truth for the response shape, shared by the
// http.HandlerFunc and the v1 API operation.
func BuildResponse(check bool, disabled bool) UpdateCheckResult {
	versionInfo := types.GetVersionInfo()

	resp := UpdateCheckResult{
		Current: CurrentVersionResponse{
			Version:   versionInfo.Version,
			Commit:    versionInfo.Commit,
			BuildTime: versionInfo.BuildTime,
			Dirty:     versionInfo.Dirty,
		},
	}

	updateCheckRequests.WithLabelValues(strconv.FormatBool(check)).Inc()

	// Only perform the remote check when requested and not disabled.
	if !check || disabled {
		return resp
	}

	// Try release version comparison for release builds.
	localVer, err := parseVersion(versionInfo.Version)
	if err == nil && !versionInfo.Dirty {
		resp = tryReleaseComparison(resp, localVer)
		if resp.Remote != nil || resp.Error != "" {
			// Release comparison produced a result (success or error).
			return resp
		}
		// Fall through to commit comparison (no releases published yet).
	}

	// Commit comparison for dev builds or when no release exists.
	resp = tryCommitComparison(resp, versionInfo.Commit)

	return resp
}

// tryReleaseComparison attempts to compare against the latest release tag.
// Returns the response with Remote set on success, Error set on fetch failure,
// or unmodified when no release exists.
func tryReleaseComparison(resp UpdateCheckResult, localVer semver) UpdateCheckResult {
	remoteTag, err := fetchRemoteLatestRelease()
	if err != nil {
		log.Warn().
			Err(err).
			Str("url", releaseAPIURL()).
			Msg("update-check: failed to fetch latest release")

		updateCheckRemoteFailures.WithLabelValues("release_fetch").Inc()

		resp.Error = "failed to check remote release: " + err.Error()

		return resp
	}

	remoteVer, err := parseVersion(remoteTag)
	if err != nil {
		// Remote tag is not a valid version — fall through to commit comparison.
		log.Warn().
			Str("tag", remoteTag).
			Msg("update-check: remote release tag is not a valid version, falling back to commit comparison")

		return resp
	}

	updateAvailable := remoteVer.GreaterThan(localVer)
	resp.UpdateAvailable = &updateAvailable
	resp.Remote = &RemoteVersionResponse{
		Commit:  remoteTag,
		Version: strings.TrimPrefix(remoteTag, "v"),
		URL:     repoURL(),
	}

	return resp
}

// tryCommitComparison compares against the latest commit on main.
func tryCommitComparison(resp UpdateCheckResult, localCommit string) UpdateCheckResult {
	remoteCommit, err := fetchRemoteShortCommit()
	if err != nil {
		log.Warn().
			Err(err).
			Str("url", commitAPIURL()).
			Msg("update-check: failed to fetch remote commit")

		updateCheckRemoteFailures.WithLabelValues("commit_fetch").Inc()

		resp.Error = "failed to check remote: " + err.Error()

		return resp
	}

	updateAvailable := remoteCommit != "" && remoteCommit != shortCommit(localCommit)
	resp.UpdateAvailable = &updateAvailable
	resp.Remote = &RemoteVersionResponse{
		Commit: remoteCommit,
		URL:    repoURL(),
	}

	return resp
}

// Handler returns an http.HandlerFunc that serves the update check endpoint.
//
// It is kept for direct mounting and tests; the v1 API registers the endpoint
// as a Huma operation backed by BuildResponse so it appears in the OpenAPI
// document. The handler:
//  1. Always returns the running binary's version info.
//  2. When ?check=true is set, fetches the remote latest release/commit
//     and compares it with the embedded version.
//
// Handler has no config access, so it always enables the remote check; the v1
// API operation passes the server's DisableUpdateCheck instead.
func Handler() http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		writer.Header().Set("Content-Type", "application/json")

		check := strings.EqualFold(req.URL.Query().Get("check"), "true")
		resp := BuildResponse(check, false)

		writer.WriteHeader(http.StatusOK)

		err := json.NewEncoder(writer).Encode(resp)
		if err != nil {
			log.Error().
				Caller().
				Err(err).
				Msg("update-check: failed to write response")
		}
	}
}

// ResetCache clears all cached remote check results. Exposed for testing.
func ResetCache() {
	cacheMu.Lock()
	defer cacheMu.Unlock()

	clear(cache)
}

// cachedFetch returns the cached value for key if fresh, otherwise calls fetch,
// stores the result, and returns it. The mutex is held during fetch so
// concurrent requests share a single remote call (singleflight behavior).
func cachedFetch(key string, ttl time.Duration, fetch func() (string, error)) (string, error) {
	cacheMu.Lock()
	defer cacheMu.Unlock()

	entry, ok := cache[key]
	if ok && time.Since(entry.fetchedAt) < ttl {
		updateCheckCacheHits.Inc()
		return entry.value, nil
	}

	updateCheckCacheMisses.Inc()

	value, err := fetch()
	if err != nil {
		return "", err
	}

	cache[key] = cacheEntry{value: value, fetchedAt: time.Now()}

	return value, nil
}

// fetchRemoteShortCommit fetches the short commit hash of the latest commit
// on the remote fork's main branch via the GitHub API. Uses caching.
func fetchRemoteShortCommit() (string, error) {
	return cachedFetch("commit", cacheTTL, func() (string, error) {
		client := &http.Client{
			Timeout: httpTimeout,
		}

		resp, err := client.Get(commitAPIURL()) //nolint:noctx // simple GET with timeout is sufficient
		if err != nil {
			return "", fmt.Errorf("github api request failed: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("%w: %d", errGitHubStatus, resp.StatusCode)
		}

		var data struct {
			SHA string `json:"sha"`
		}

		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return "", fmt.Errorf("failed to decode github api response: %w", err)
		}

		if data.SHA == "" {
			return "", errGitHubEmptySHA
		}

		return shortCommit(data.SHA), nil
	})
}

// fetchRemoteLatestRelease fetches the latest release tag from the remote repo
// via the GitHub API. Uses caching.
func fetchRemoteLatestRelease() (string, error) {
	return cachedFetch("release", cacheTTL, func() (string, error) {
		client := &http.Client{
			Timeout: httpTimeout,
		}

		resp, err := client.Get(releaseAPIURL()) //nolint:noctx // simple GET with timeout is sufficient
		if err != nil {
			return "", fmt.Errorf("github api request failed: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("%w: %d", errGitHubStatus, resp.StatusCode)
		}

		var data struct {
			TagName string `json:"tag_name"`
		}

		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return "", fmt.Errorf("failed to decode github api response: %w", err)
		}

		if data.TagName == "" {
			return "", errGitHubEmptyTag
		}

		return data.TagName, nil
	})
}

// shortCommit returns the first 7 characters of a commit hash, or the
// full string if it is already shorter than 7 characters.
func shortCommit(commit string) string {
	if len(commit) > 7 {
		return commit[:7]
	}

	return commit
}

// BuildInfo returns the vcs.revision and vcs.time from the build info,
// falling back to the VersionInfo's commit if build info is unavailable.
// This is used for embedding additional build metadata.
func BuildInfo() (string, string) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return unknownVersion, unknownVersion
	}

	var revision, time string

	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			revision = setting.Value
		case "vcs.time":
			time = setting.Value
		}
	}

	if revision == "" {
		revision = unknownVersion
	}

	if time == "" {
		time = unknownVersion
	}

	return revision, time
}

// parseVersion parses a semantic version string (e.g., "v0.29.9-arsydoni4326-alt"
// or "0.29.9") into a semver struct. Returns an error for non-version strings.
func parseVersion(s string) (semver, error) {
	if isDevVersion(s) {
		return semver{}, fmt.Errorf("%w: %q", errNotAVersion, s)
	}

	s = strings.TrimPrefix(s, "v")

	var pre string
	if idx := strings.IndexAny(s, "-+"); idx != -1 {
		pre = s[idx:]
		s = s[:idx]
	}

	parts := strings.Split(s, ".")
	if len(parts) != 3 {
		return semver{}, fmt.Errorf("%w: %q", errVersionParts, s)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return semver{}, fmt.Errorf("%w: %q", errMajorNotNumber, parts[0])
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return semver{}, fmt.Errorf("%w: %q", errMinorNotNumber, parts[1])
	}

	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return semver{}, fmt.Errorf("%w: %q", errPatchNotNumber, parts[2])
	}

	return semver{major: major, minor: minor, patch: patch, preRelease: pre}, nil
}

// pseudoVersionRe matches Go module pseudo-versions
// (vX.Y.Z-<14-digit timestamp>-<12-char hash>), which identify untagged
// main-sha builds and must be treated as dev builds.
var pseudoVersionRe = regexp.MustCompile(`^v?[0-9]+\.[0-9]+\.[0-9]+-[0-9]{14}-[0-9a-f]{12}$`)

// isDevVersion reports whether a version string represents a development build.
func isDevVersion(s string) bool {
	if s == "" || s == "dev" || s == unknownVersion || s == "(devel)" {
		return true
	}

	return pseudoVersionRe.MatchString(s)
}

// GreaterThan returns true if v > other, following semver precedence:
// major > minor > patch, and no pre-release > pre-release.
func (v semver) GreaterThan(other semver) bool {
	switch {
	case v.major != other.major:
		return v.major > other.major
	case v.minor != other.minor:
		return v.minor > other.minor
	case v.patch != other.patch:
		return v.patch > other.patch
	default:
		// Same major.minor.patch — compare pre-release.
		if v.preRelease == other.preRelease {
			return false
		}
		// No pre-release > pre-release (e.g., 0.29.9 > 0.29.9-arsydoni4326-alt).
		if v.preRelease == "" {
			return true
		}

		if other.preRelease == "" {
			return false
		}
		// Both have pre-release — lexical comparison.
		return v.preRelease > other.preRelease
	}
}
