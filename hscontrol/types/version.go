package types

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
)

type GoInfo struct {
	Version string `json:"version"`
	OS      string `json:"os"`
	Arch    string `json:"arch"`
}

type VersionInfo struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildTime string `json:"buildTime"`
	Go        GoInfo `json:"go"`
	Dirty     bool   `json:"dirty"`
}

func (v *VersionInfo) String() string {
	var sb strings.Builder

	version := v.Version
	if v.Dirty && !strings.Contains(version, "dirty") {
		version += "-dirty"
	}

	fmt.Fprintf(&sb, "headscale version %s\n", version)
	fmt.Fprintf(&sb, "commit: %s\n", v.Commit)
	fmt.Fprintf(&sb, "build time: %s\n", v.BuildTime)
	fmt.Fprintf(&sb, "built with: %s %s/%s\n", v.Go.Version, v.Go.OS, v.Go.Arch)

	return sb.String()
}

// Build-time injected variables.
//
// These are set via -ldflags -X by the Docker/CI pipeline:
//
//	-X 'github.com/arsydoni4326-alt/headscale/hscontrol/types.Version=...'
//	-X 'github.com/arsydoni4326-alt/headscale/hscontrol/types.Commit=...'
//	-X 'github.com/arsydoni4326-alt/headscale/hscontrol/types.BuildDate=...'
//
// They take precedence over the VCS info embedded by the Go toolchain, which
// is unavailable in Docker builds because .git is excluded from the build
// context.
var (
	Version   string
	Commit    string
	BuildDate string
)

var buildInfo = sync.OnceValues(debug.ReadBuildInfo)

var GetVersionInfo = sync.OnceValue(func() *VersionInfo {
	bi, ok := buildInfo()
	return buildVersionInfo(Version, Commit, BuildDate, bi, ok)
})

// buildVersionInfo assembles the VersionInfo from the Go toolchain's embedded
// build info (if any) and the build-time injected values, which take
// precedence. Extracted as a pure function so the precedence rules are
// testable.
func buildVersionInfo(
	injectedVersion, injectedCommit, injectedBuildDate string,
	bi *debug.BuildInfo,
	ok bool,
) *VersionInfo {
	info := &VersionInfo{
		Version:   "dev",
		Commit:    "unknown",
		BuildTime: "unknown",
		Go: GoInfo{
			Version: runtime.Version(),
			OS:      runtime.GOOS,
			Arch:    runtime.GOARCH,
		},
	}

	if ok && bi != nil {
		// Extract version from module path or main version
		if bi.Main.Version != "" && bi.Main.Version != "(devel)" {
			info.Version = bi.Main.Version
		}

		// Extract build settings
		for _, setting := range bi.Settings {
			switch setting.Key {
			case "vcs.revision":
				info.Commit = setting.Value
			case "vcs.modified":
				info.Dirty = setting.Value == "true"
			case "vcs.time":
				info.BuildTime = setting.Value
			}
		}
	}

	// Build-time injected values take precedence over the toolchain's
	// embedded VCS info.
	if injectedVersion != "" {
		info.Version = injectedVersion
	}

	if injectedCommit != "" {
		info.Commit = injectedCommit
	}

	if injectedBuildDate != "" {
		info.BuildTime = injectedBuildDate
	}

	return info
}
