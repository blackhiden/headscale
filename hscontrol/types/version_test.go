package types

import (
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildVersionInfoDefaults(t *testing.T) {
	info := buildVersionInfo("", "", "", nil, false)

	assert.Equal(t, "dev", info.Version)
	assert.Equal(t, "unknown", info.Commit)
	assert.Equal(t, "unknown", info.BuildTime)
	assert.False(t, info.Dirty)
}

func TestBuildVersionInfoFromBuildInfo(t *testing.T) {
	bi := &debug.BuildInfo{
		Main: debug.Module{
			Path:    "github.com/arsydoni4326-alt/headscale",
			Version: "v0.29.11-arsydoni4326-alt",
		},
		Settings: []debug.BuildSetting{
			{Key: "vcs.revision", Value: "7f5e0a3c1234567890abcdef1234567890abcdef12"},
			{Key: "vcs.modified", Value: "false"},
			{Key: "vcs.time", Value: "2026-09-26T10:00:00Z"},
		},
	}

	info := buildVersionInfo("", "", "", bi, true)

	assert.Equal(t, "v0.29.11-arsydoni4326-alt", info.Version)
	assert.Equal(t, "7f5e0a3c1234567890abcdef1234567890abcdef12", info.Commit)
	assert.Equal(t, "2026-09-26T10:00:00Z", info.BuildTime)
	assert.False(t, info.Dirty)
}

func TestBuildVersionInfoInjectedTakesPrecedence(t *testing.T) {
	bi := &debug.BuildInfo{
		Main: debug.Module{Path: "github.com/arsydoni4326-alt/headscale", Version: "(devel)"},
		Settings: []debug.BuildSetting{
			{Key: "vcs.revision", Value: "oldcommit"},
			{Key: "vcs.modified", Value: "true"},
			{Key: "vcs.time", Value: "2020-01-01T00:00:00Z"},
		},
	}

	info := buildVersionInfo(
		"v0.29.11-arsydoni4326-alt",
		"7f5e0a3c",
		"2026-09-27T12:00:00Z",
		bi,
		true,
	)

	assert.Equal(t, "v0.29.11-arsydoni4326-alt", info.Version)
	assert.Equal(t, "7f5e0a3c", info.Commit)
	assert.Equal(t, "2026-09-27T12:00:00Z", info.BuildTime)
	// Dirty is not injected; it still comes from the build info.
	assert.True(t, info.Dirty)
}

func TestBuildVersionInfoInjectedEmptyKeepsBuildInfo(t *testing.T) {
	bi := &debug.BuildInfo{
		Main: debug.Module{
			Path:    "github.com/arsydoni4326-alt/headscale",
			Version: "v0.29.10-arsydoni4326-alt",
		},
		Settings: []debug.BuildSetting{
			{Key: "vcs.revision", Value: "abcdef1"},
			{Key: "vcs.time", Value: "2026-09-25T00:00:00Z"},
		},
	}

	// Only version is injected; commit and build time come from build info.
	info := buildVersionInfo("v0.29.11-arsydoni4326-alt", "", "", bi, true)

	assert.Equal(t, "v0.29.11-arsydoni4326-alt", info.Version)
	assert.Equal(t, "abcdef1", info.Commit)
	assert.Equal(t, "2026-09-25T00:00:00Z", info.BuildTime)
}
