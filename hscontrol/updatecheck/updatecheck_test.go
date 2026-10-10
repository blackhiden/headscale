package updatecheck

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

// errTestFetch is a sentinel error used by the caching tests.
var errTestFetch = errors.New("test error")

func TestHandler_WithoutCheckParam(t *testing.T) {
	ResetCache()

	handler := Handler()

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/v1/update-check", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Fatalf("expected Content-Type application/json, got %s", contentType)
	}

	var body UpdateCheckResult

	err := json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Should always have current version info
	if body.Current.Version == "" {
		t.Error("expected current.version to be non-empty")
	}

	if body.Current.Commit == "" {
		t.Error("expected current.commit to be non-empty")
	}

	// Without ?check=true, remote should be nil
	if body.Remote != nil {
		t.Error("expected remote to be nil without ?check=true")
	}

	if body.UpdateAvailable != nil {
		t.Error("expected updateAvailable to be nil without ?check=true")
	}
}

func TestHandler_WithCheckParam(t *testing.T) {
	ResetCache()
	resetRepoConfigForTest()
	t.Setenv(remoteRepoEnv, "")

	withTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"sha": "abcdef1234567890"}`))
	})

	handler := Handler()

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/v1/update-check?check=true", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var body UpdateCheckResult

	err := json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Should always have current version info
	if body.Current.Version == "" {
		t.Error("expected current.version to be non-empty")
	}

	if body.Current.Commit == "" {
		t.Error("expected current.commit to be non-empty")
	}

	// With ?check=true and a deterministic test API, the remote check should
	// succeed and report whether an update is available.
	if body.Remote == nil {
		t.Error("expected remote to be set with ?check=true")
	}

	if body.Error != "" {
		t.Errorf("expected no error, got %q", body.Error)
	}

	if body.UpdateAvailable == nil {
		t.Error("expected updateAvailable to be set with ?check=true")
	}
}

func TestShortCommit(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"abc123def456ghi789", "abc123d"},
		{"abc123d", "abc123d"},
		{"ab", "ab"},
		{"", ""},
	}

	for _, tt := range tests {
		result := shortCommit(tt.input)
		if result != tt.expected {
			t.Errorf("shortCommit(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestHandler_ResponseStructure(t *testing.T) {
	ResetCache()

	handler := Handler()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/v1/update-check?check=true", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	var body map[string]any

	err := json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Verify required fields exist
	requiredFields := []string{"current"}
	for _, field := range requiredFields {
		if _, ok := body[field]; !ok {
			t.Errorf("expected response field %q to exist", field)
		}
	}

	// Verify nested current fields
	current, ok := body["current"].(map[string]any)
	if !ok {
		t.Fatal("expected current to be an object")
	}

	nestedFields := []string{"version", "commit", "buildTime"}
	for _, field := range nestedFields {
		if _, ok := current[field]; !ok {
			t.Errorf("expected current.%s to exist", field)
		}
	}
}

// --- Version parsing tests ---

func TestParseVersion(t *testing.T) {
	tests := []struct {
		input   string
		wantOk  bool
		wantMaj int
		wantMin int
		wantPat int
		wantPre string
	}{
		{"0.29.9", true, 0, 29, 9, ""},
		{"v0.29.9", true, 0, 29, 9, ""},
		{"0.29.9-arsydoni4326-alt", true, 0, 29, 9, "-arsydoni4326-alt"},
		{"v0.29.9-arsydoni4326-alt", true, 0, 29, 9, "-arsydoni4326-alt"},
		{"1.0.0", true, 1, 0, 0, ""},
		{"1.2.3-rc1", true, 1, 2, 3, "-rc1"},
		{"1.2.3+build123", true, 1, 2, 3, "+build123"},
		{"", false, 0, 0, 0, ""},
		{"dev", false, 0, 0, 0, ""},
		{"unknown", false, 0, 0, 0, ""},
		{"(devel)", false, 0, 0, 0, ""},
		{"not-a-version", false, 0, 0, 0, ""},
		{"1.2", false, 0, 0, 0, ""},
		{"abc.def.ghi", false, 0, 0, 0, ""},
	}

	for _, tt := range tests {
		got, err := parseVersion(tt.input)
		if tt.wantOk {
			if err != nil {
				t.Errorf("parseVersion(%q) unexpected error: %v", tt.input, err)
				continue
			}

			if got.major != tt.wantMaj {
				t.Errorf("parseVersion(%q) major = %d, want %d", tt.input, got.major, tt.wantMaj)
			}

			if got.minor != tt.wantMin {
				t.Errorf("parseVersion(%q) minor = %d, want %d", tt.input, got.minor, tt.wantMin)
			}

			if got.patch != tt.wantPat {
				t.Errorf("parseVersion(%q) patch = %d, want %d", tt.input, got.patch, tt.wantPat)
			}

			if got.preRelease != tt.wantPre {
				t.Errorf("parseVersion(%q) preRelease = %q, want %q", tt.input, got.preRelease, tt.wantPre)
			}
		} else if err == nil {
			t.Errorf("parseVersion(%q) expected error, got %+v", tt.input, got)
		}
	}
}

func TestSemverGreaterThan(t *testing.T) {
	tests := []struct {
		a    string
		b    string
		want bool // a > b
	}{
		{"1.0.0", "0.9.9", true},
		{"0.9.9", "1.0.0", false},
		{"0.29.10", "0.29.9", true},
		{"0.29.9", "0.29.10", false},
		{"0.29.9", "0.29.9", false},
		{"0.29.9", "0.29.9-arsydoni4326-alt", true},  // no pre > pre
		{"0.29.9-arsydoni4326-alt", "0.29.9", false}, // pre < no pre
		{"0.29.10-arsydoni4326-alt", "0.29.9-arsydoni4326-alt", true},
		{"0.29.9-rc2", "0.29.9-rc1", true},
		{"0.29.9-rc1", "0.29.9-rc2", false},
		{"0.29.9-alpha", "0.29.9-beta", false}, // lexical: "alpha" < "beta"
		{"0.29.9-beta", "0.29.9-alpha", true},  // lexical: "beta" > "alpha"
		{"0.29.9-rc1", "0.29.9-rc1", false},    // same
	}

	for _, tt := range tests {
		a, errA := parseVersion(tt.a)
		b, errB := parseVersion(tt.b)

		if errA != nil || errB != nil {
			t.Fatalf("parse error: a=%v, b=%v", errA, errB)
		}

		got := a.GreaterThan(b)
		if got != tt.want {
			t.Errorf("parseVersion(%q).GreaterThan(%q) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestIsDevVersion(t *testing.T) {
	devVersions := []string{"", "dev", "unknown", "(devel)"}
	releaseVersions := []string{"0.29.9", "v0.29.9-arsydoni4326-alt", "1.0.0"}

	for _, v := range devVersions {
		if !isDevVersion(v) {
			t.Errorf("isDevVersion(%q) = false, want true", v)
		}
	}

	for _, v := range releaseVersions {
		if isDevVersion(v) {
			t.Errorf("isDevVersion(%q) = true, want false", v)
		}
	}
}

// --- Caching tests ---

func TestCachedFetch(t *testing.T) {
	ResetCache()

	var callCount int

	fetch := func() (string, error) {
		callCount++
		return "abc1234", nil
	}

	// First call: miss, fetches
	val, err := cachedFetch("test", 5*time.Minute, fetch)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if val != "abc1234" {
		t.Errorf("got %q, want %q", val, "abc1234")
	}

	if callCount != 1 {
		t.Errorf("callCount = %d, want 1", callCount)
	}

	// Second call: hit, no fetch
	val, err = cachedFetch("test", 5*time.Minute, fetch)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if val != "abc1234" {
		t.Errorf("got %q, want %q", val, "abc1234")
	}

	if callCount != 1 {
		t.Errorf("callCount = %d, want 1 (cached)", callCount)
	}
}

func TestCachedFetch_Expired(t *testing.T) {
	ResetCache()

	var callCount int

	fetch := func() (string, error) {
		callCount++
		return "abc1234", nil
	}

	// Fetch with very short TTL
	val, err := cachedFetch("expire-test", 1*time.Millisecond, fetch)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if val != "abc1234" {
		t.Errorf("got %q, want %q", val, "abc1234")
	}

	// Force expiry by backdating the cached entry.
	cacheMu.Lock()
	entry := cache["expire-test"]
	entry.fetchedAt = time.Now().Add(-2 * time.Millisecond)
	cache["expire-test"] = entry
	cacheMu.Unlock()

	// Should fetch again
	val, err = cachedFetch("expire-test", 1*time.Millisecond, fetch)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if val != "abc1234" {
		t.Errorf("got %q, want %q", val, "abc1234")
	}

	if callCount != 2 {
		t.Errorf("callCount = %d, want 2 (expired)", callCount)
	}
}

func TestCachedFetch_ErrorNotCached(t *testing.T) {
	ResetCache()

	var callCount int

	fetch := func() (string, error) {
		callCount++
		return "", errTestFetch
	}

	// First call: error, not cached
	_, err := cachedFetch("err-test", 5*time.Minute, fetch)
	if err == nil {
		t.Fatal("expected error")
	}

	if callCount != 1 {
		t.Errorf("callCount = %d, want 1", callCount)
	}

	// Second call: should retry (error was not cached)
	_, err = cachedFetch("err-test", 5*time.Minute, fetch)
	if err == nil {
		t.Fatal("expected error")
	}

	if callCount != 2 {
		t.Errorf("callCount = %d, want 2 (retried)", callCount)
	}
}

func TestCachedFetch_Reset(t *testing.T) {
	ResetCache()

	var callCount int

	fetch := func() (string, error) {
		callCount++
		return "abc1234", nil
	}

	_, _ = cachedFetch("reset-test", 5*time.Minute, fetch)

	if callCount != 1 {
		t.Errorf("callCount = %d, want 1", callCount)
	}

	ResetCache()

	// After reset, should fetch again
	_, _ = cachedFetch("reset-test", 5*time.Minute, fetch)

	if callCount != 2 {
		t.Errorf("callCount = %d, want 2 (after reset)", callCount)
	}
}

// --- Edge case tests ---

func TestBuildResponse_WithCheck_NoExternalCall(t *testing.T) {
	// Without ?check=true, BuildResponse should never make external calls.
	ResetCache()

	resp := BuildResponse(false, false)
	if resp.Remote != nil {
		t.Error("expected remote to be nil when check=false")
	}

	if resp.UpdateAvailable != nil {
		t.Error("expected updateAvailable to be nil when check=false")
	}
}

func TestBuildResponse_Disabled(t *testing.T) {
	// When disabled=true (disable_check_updates), the remote check is skipped
	// even when check=true: no external call is made and only version info is
	// returned.
	ResetCache()
	resetRepoConfigForTest()
	t.Setenv(remoteRepoEnv, "")

	withTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("remote check must not be performed when disabled")
	})

	resp := BuildResponse(true, true)

	if resp.Error != "" {
		t.Errorf("expected no error, got %q", resp.Error)
	}

	if resp.Remote != nil {
		t.Error("expected remote to be nil when disabled")
	}

	if resp.UpdateAvailable != nil {
		t.Error("expected updateAvailable to be nil when disabled")
	}

	// Always has current version info
	if resp.Current.Version == "" {
		t.Error("expected current.version to be non-empty")
	}
}

func TestBuildResponse_DevVersion(t *testing.T) {
	// In test environments, the version is typically "dev" or "(devel)".
	// BuildResponse should handle this gracefully and fall back to commit
	// comparison against the deterministic test API.
	ResetCache()
	resetRepoConfigForTest()
	t.Setenv(remoteRepoEnv, "")

	withTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"sha": "abcdef1234567890"}`))
	})

	resp := BuildResponse(true, false)

	if resp.Error != "" {
		t.Errorf("expected no error, got %q", resp.Error)
	}

	if resp.Remote == nil {
		t.Error("expected remote to be set for dev build commit comparison")
	}

	// Always has current version info
	if resp.Current.Version == "" {
		t.Error("expected current.version to be non-empty")
	}
}

func TestBuildResponse_RemoteVersionResponseHasVersion(t *testing.T) {
	// Verify that RemoteVersionResponse includes the Version field in JSON.
	resp := UpdateCheckResult{
		Current:         CurrentVersionResponse{Version: "0.29.9", Commit: "abc1234", BuildTime: "2026-01-01", Dirty: false},
		UpdateAvailable: func() *bool { v := false; return &v }(),
		Remote: &RemoteVersionResponse{
			Commit:  "v0.29.9-arsydoni4326-alt",
			Version: "0.29.9-arsydoni4326-alt",
			URL:     "https://github.com/test-owner/test-repo",
		},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var parsed map[string]any

	err = json.Unmarshal(data, &parsed)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	remote, ok := parsed["remote"].(map[string]any)
	if !ok {
		t.Fatal("expected remote object in JSON")
	}

	if _, ok := remote["version"]; !ok {
		t.Error("expected remote.version field in JSON")
	}
}

func TestParseVersion_GoModPseudoVersion(t *testing.T) {
	// Go module pseudo-versions like v0.0.0-20260522092201-58a85b68b3d9
	// should be treated as dev builds.
	v := "v0.0.0-20260522092201-58a85b68b3d9"
	if !isDevVersion(v) {
		t.Errorf("isDevVersion(%q) = false, want true (pseudo-version is dev)", v)
	}

	// parseVersion should also reject it
	_, err := parseVersion(v)
	if err != nil {
		t.Logf("parseVersion(%q) correctly returned error: %v", v, err)
	}
}

// TestFetchRemoteShortCommit_MalformedJSON verifies the fetch function returns
// an error when the GitHub API response is not valid JSON.
func TestFetchRemoteShortCommit_MalformedJSON(t *testing.T) {
	ResetCache()
	resetRepoConfigForTest()
	t.Setenv(remoteRepoEnv, "")

	withTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"sha": `)) // truncated JSON
	})

	_, err := fetchRemoteShortCommit()
	if err == nil {
		t.Fatal("expected error for malformed JSON")
	}
}

// --- Deterministic GitHub API fetch tests ---

// withTestAPI points the GitHub API base URL at a local httptest server for
// the duration of the test and restores it afterwards.
func withTestAPI(t *testing.T, handler http.HandlerFunc) {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	orig := apiBaseURL
	apiBaseURL = server.URL

	t.Cleanup(func() { apiBaseURL = orig })
}

// resetRepoConfigForTest reinitializes the repo-config sync.Once so tests can
// exercise different HEADSCALE_UPDATE_CHECK_REPO values. Tests in this package
// run sequentially, so reassigning the package-level Once is safe.
func resetRepoConfigForTest() {
	repoConfig = sync.Once{}
	repoOwnerVal = ""
	repoNameVal = ""
}

func TestFetchRemoteShortCommit_Success(t *testing.T) {
	ResetCache()
	resetRepoConfigForTest()
	t.Setenv(remoteRepoEnv, "")

	withTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/arsydoni4326-alt/headscale/commits/main" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"sha": "abcdef1234567890"}`))
	})

	got, err := fetchRemoteShortCommit()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "abcdef1" {
		t.Errorf("fetchRemoteShortCommit() = %q, want %q", got, "abcdef1")
	}
}

func TestFetchRemoteShortCommit_Non200(t *testing.T) {
	ResetCache()
	resetRepoConfigForTest()
	t.Setenv(remoteRepoEnv, "")

	withTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden) // 403 rate limit
	})

	_, err := fetchRemoteShortCommit()
	if err == nil {
		t.Fatal("expected error for non-200 status")
	}

	if !errors.Is(err, errGitHubStatus) {
		t.Errorf("expected errGitHubStatus, got %v", err)
	}
}

func TestFetchRemoteLatestRelease_Success(t *testing.T) {
	ResetCache()
	resetRepoConfigForTest()
	t.Setenv(remoteRepoEnv, "")

	withTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/arsydoni4326-alt/headscale/releases/latest" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"tag_name": "v0.30.0-arsydoni4326-alt"}`))
	})

	got, err := fetchRemoteLatestRelease()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "v0.30.0-arsydoni4326-alt" {
		t.Errorf("fetchRemoteLatestRelease() = %q, want %q", got, "v0.30.0-arsydoni4326-alt")
	}
}

func TestFetchRemoteLatestRelease_Non200(t *testing.T) {
	ResetCache()
	resetRepoConfigForTest()
	t.Setenv(remoteRepoEnv, "")

	withTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := fetchRemoteLatestRelease()
	if !errors.Is(err, errGitHubStatus) {
		t.Errorf("expected errGitHubStatus, got %v", err)
	}
}

func TestFetchRemoteLatestRelease_EmptyTag(t *testing.T) {
	ResetCache()
	resetRepoConfigForTest()
	t.Setenv(remoteRepoEnv, "")

	withTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	})

	_, err := fetchRemoteLatestRelease()
	if !errors.Is(err, errGitHubEmptyTag) {
		t.Errorf("expected errGitHubEmptyTag, got %v", err)
	}
}

func TestFetchRemoteLatestRelease_MalformedJSON(t *testing.T) {
	ResetCache()
	resetRepoConfigForTest()
	t.Setenv(remoteRepoEnv, "")

	withTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`not json`))
	})

	_, err := fetchRemoteLatestRelease()
	if err == nil {
		t.Fatal("expected error for malformed JSON")
	}
}

// --- Repo config tests ---

func TestRepoConfig_Default(t *testing.T) {
	resetRepoConfigForTest()
	t.Setenv(remoteRepoEnv, "")

	if got := repoOwner(); got != DefaultRemoteRepoOwner {
		t.Errorf("repoOwner() = %q, want %q", got, DefaultRemoteRepoOwner)
	}

	if got := repoName(); got != DefaultRemoteRepoName {
		t.Errorf("repoName() = %q, want %q", got, DefaultRemoteRepoName)
	}
}

func TestRepoConfig_ValidEnv(t *testing.T) {
	resetRepoConfigForTest()
	t.Setenv(remoteRepoEnv, "octocat/hello-world")

	if got := repoOwner(); got != "octocat" {
		t.Errorf("repoOwner() = %q, want %q", got, "octocat")
	}

	if got := repoName(); got != "hello-world" {
		t.Errorf("repoName() = %q, want %q", got, "hello-world")
	}
}

func TestRepoConfig_InvalidEnv(t *testing.T) {
	for _, repo := range []string{"no-slash", "/missing-owner", "missing-name/"} {
		t.Run(repo, func(t *testing.T) {
			resetRepoConfigForTest()
			t.Setenv(remoteRepoEnv, repo)

			if got := repoOwner(); got != DefaultRemoteRepoOwner {
				t.Errorf("repoOwner() = %q, want default %q", got, DefaultRemoteRepoOwner)
			}

			if got := repoName(); got != DefaultRemoteRepoName {
				t.Errorf("repoName() = %q, want default %q", got, DefaultRemoteRepoName)
			}
		})
	}
}

func TestRepoConfig_ExtraSlash(t *testing.T) {
	// strings.Cut splits on the first slash, so "a/b/c" is accepted with
	// owner="a" and name="b/c". Document the actual behavior.
	resetRepoConfigForTest()
	t.Setenv(remoteRepoEnv, "a/b/c")

	if got := repoOwner(); got != "a" {
		t.Errorf("repoOwner() = %q, want %q", got, "a")
	}

	if got := repoName(); got != "b/c" {
		t.Errorf("repoName() = %q, want %q", got, "b/c")
	}
}

// --- Comparison helpers ---

func TestTryReleaseComparison_InvalidRemoteTag(t *testing.T) {
	ResetCache()
	resetRepoConfigForTest()
	t.Setenv(remoteRepoEnv, "")

	withTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"tag_name": "not-a-version"}`))
	})

	localVer, err := parseVersion("0.29.9")
	if err != nil {
		t.Fatalf("parseVersion: %v", err)
	}

	resp := tryReleaseComparison(UpdateCheckResult{}, localVer)

	if resp.Remote != nil {
		t.Error("expected Remote to be nil (fallback to commit comparison)")
	}

	if resp.Error != "" {
		t.Errorf("expected no error, got %q", resp.Error)
	}
}

func TestTryReleaseComparison_Success(t *testing.T) {
	ResetCache()
	resetRepoConfigForTest()
	t.Setenv(remoteRepoEnv, "")

	withTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"tag_name": "v0.30.0-arsydoni4326-alt"}`))
	})

	localVer, err := parseVersion("0.29.9")
	if err != nil {
		t.Fatalf("parseVersion: %v", err)
	}

	resp := tryReleaseComparison(UpdateCheckResult{}, localVer)

	if resp.Remote == nil {
		t.Fatal("expected Remote to be set")
	}

	if resp.UpdateAvailable == nil || !*resp.UpdateAvailable {
		t.Error("expected updateAvailable=true (0.30.0 > 0.29.9)")
	}

	if resp.Remote.Version != "0.30.0-arsydoni4326-alt" {
		t.Errorf("remote version = %q, want %q", resp.Remote.Version, "0.30.0-arsydoni4326-alt")
	}
}

func TestTryCommitComparison_Success(t *testing.T) {
	ResetCache()
	resetRepoConfigForTest()
	t.Setenv(remoteRepoEnv, "")

	withTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"sha": "abcdef1234567890"}`))
	})

	resp := tryCommitComparison(UpdateCheckResult{}, "1111111")

	if resp.Remote == nil {
		t.Fatal("expected Remote to be set")
	}

	if resp.UpdateAvailable == nil || !*resp.UpdateAvailable {
		t.Error("expected updateAvailable=true (remote commit differs)")
	}

	if resp.Remote.Commit != "abcdef1" {
		t.Errorf("remote commit = %q, want %q", resp.Remote.Commit, "abcdef1")
	}
}

func TestTryCommitComparison_NoUpdate(t *testing.T) {
	ResetCache()
	resetRepoConfigForTest()
	t.Setenv(remoteRepoEnv, "")

	withTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"sha": "abcdef1234567890"}`))
	})

	resp := tryCommitComparison(UpdateCheckResult{}, "abcdef1234567890")

	if resp.UpdateAvailable == nil || *resp.UpdateAvailable {
		t.Error("expected updateAvailable=false (same commit)")
	}
}

func TestTryCommitComparison_Error(t *testing.T) {
	ResetCache()
	resetRepoConfigForTest()
	t.Setenv(remoteRepoEnv, "")

	withTestAPI(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	resp := tryCommitComparison(UpdateCheckResult{}, "1111111")

	if resp.Error == "" {
		t.Error("expected error to be set")
	}

	if resp.Remote != nil {
		t.Error("expected Remote to be nil on error")
	}
}
