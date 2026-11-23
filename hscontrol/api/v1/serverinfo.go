package apiv1

import (
	"context"
	"fmt"
	"net/http"
	"net/netip"
	"os"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/blackhiden/headscale/hscontrol/types"
	"github.com/danielgtaylor/huma/v2"
)

// CPUInfo describes the server's CPU configuration.
type CPUInfo struct {
	Count int `json:"count"`
}

// MemoryInfo describes the server's memory configuration.
type MemoryInfo struct {
	TotalBytes uint64 `json:"totalBytes"`
}

// StorageInfo describes the server's storage configuration.
type StorageInfo struct {
	TotalBytes uint64 `json:"totalBytes"`
}

// DeploymentInfo describes how Headscale is deployed.
type DeploymentInfo struct {
	Mode          string `json:"mode"`
	DockerVersion string `json:"dockerVersion,omitempty"`
}

// TailnetInfo describes the tailnet network configuration.
type TailnetInfo struct {
	IPRangeV4      string `json:"ipRangeV4"`
	IPRangeV6      string `json:"ipRangeV6"`
	MagicDNSDomain string `json:"magicDnsDomain"`
	DERPRegions    int    `json:"derpRegions"`
}

// ServerInfoResponseBody is the response for GET /api/v1/server-info.
type ServerInfoResponseBody struct {
	Version    string         `json:"version"`
	Commit     string         `json:"commit"`
	BuildTime  string         `json:"buildTime"`
	Go         types.GoInfo   `json:"go"`
	OS         string         `json:"os"`
	Arch       string         `json:"arch"`
	CPU        CPUInfo        `json:"cpu"`
	Memory     MemoryInfo     `json:"memory"`
	Storage    StorageInfo    `json:"storage"`
	Deployment DeploymentInfo `json:"deployment"`
	PublicIP   string         `json:"publicIp"`
	Tailnet    TailnetInfo    `json:"tailnet"`
}

type serverInfoOutput struct {
	Body ServerInfoResponseBody
}

func init() {
	registrations = append(registrations, registerServerInfo)
}

func registerServerInfo(api huma.API, b Backend) {
	huma.Register(api, huma.Operation{
		OperationID: "getServerInfo",
		Method:      http.MethodGet,
		Path:        "/api/v1/server-info",
		Summary:     "Get server technical info",
		Description: "Returns non-sensitive technical information about the Headscale server.",
		Tags:        []string{"Server"},
		Security:    bearerAuth,
	}, func(ctx context.Context, _ *struct{}) (*serverInfoOutput, error) {
		versionInfo := types.GetVersionInfo()

		derpMap := b.State.DERPMap()

		var derpRegions int
		if derpMap.Valid() {
			derpRegions = derpMap.Regions().Len()
		}

		body := ServerInfoResponseBody{
			Version:   versionInfo.Version,
			Commit:    versionInfo.Commit,
			BuildTime: versionInfo.BuildTime,
			Go:        versionInfo.Go,
			OS:        runtime.GOOS,
			Arch:      runtime.GOARCH,
			CPU: CPUInfo{
				Count: runtime.NumCPU(),
			},
			Memory: MemoryInfo{
				TotalBytes: readMemTotal(),
			},
			Storage: StorageInfo{
				TotalBytes: readStorageTotal(),
			},
			Deployment: DeploymentInfo{
				Mode:          detectDeploymentMode(),
				DockerVersion: detectDockerVersion(),
			},
			PublicIP: detectPublicIP(b),
			Tailnet: TailnetInfo{
				IPRangeV4:      prefixString(b.Cfg.Prefixes.V4),
				IPRangeV6:      prefixString(b.Cfg.Prefixes.V6),
				MagicDNSDomain: b.Cfg.DNSConfig.BaseDomain,
				DERPRegions:    derpRegions,
			},
		}

		return &serverInfoOutput{Body: body}, nil
	})
}

// prefixString formats a netip.Prefix as a string, returning "" for a nil prefix.
func prefixString(p *netip.Prefix) string {
	if p == nil {
		return ""
	}

	return p.String()
}

// readMemTotal reads total system memory from /proc/meminfo on Linux.
// Returns 0 on non-Linux or when the file cannot be read.
func readMemTotal() uint64 {
	if runtime.GOOS != "linux" {
		return 0
	}

	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return 0
	}

	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "MemTotal:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				var kb uint64
				if _, err := fmt.Sscanf(fields[1], "%d", &kb); err == nil {
					return kb * 1024
				}
			}
		}
	}

	return 0
}

// readStorageTotal reads total storage for the working directory.
// Returns 0 on non-Linux or when the statfs call fails.
func readStorageTotal() uint64 {
	if runtime.GOOS != "linux" {
		return 0
	}

	var stat syscall.Statfs_t
	if err := syscall.Statfs(".", &stat); err != nil {
		return 0
	}

	return stat.Blocks * uint64(stat.Bsize)
}

// detectDeploymentMode returns "docker" when running in a container, "native" otherwise.
func detectDeploymentMode() string {
	if isInDocker() {
		return "docker"
	}

	return "native"
}

// isInDocker checks for common Docker/container indicators.
func isInDocker() bool {
	// Check for /.dockerenv
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	// Check for /run/.containerenv
	if _, err := os.Stat("/run/.containerenv"); err == nil {
		return true
	}

	// Check /proc/1/cgroup for docker/containerd entries
	data, err := os.ReadFile("/proc/1/cgroup")
	if err == nil {
		content := string(data)
		if strings.Contains(content, "docker") || strings.Contains(content, "containerd") {
			return true
		}
	}

	return false
}

// detectDockerVersion attempts to determine the Docker version when running in Docker.
// Returns "unknown" when not determinable.
func detectDockerVersion() string {
	if !isInDocker() {
		return ""
	}

	// Docker version is not directly available from inside the container,
	// so report "unknown" unless we can determine it.
	return "unknown"
}

// detectPublicIP attempts to determine the server's public IP via an outbound
// HTTP request to a public IP echo service. Falls back to the configured
// server_url hostname. Non-blocking: never fails the endpoint.
func detectPublicIP(b Backend) string {
	// Try the public IP echo service with a short timeout.
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("https://api.ipify.org")
	if err == nil {
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			buf := make([]byte, 64)
			n, _ := resp.Body.Read(buf)
			ip := strings.TrimSpace(string(buf[:n]))
			if ip != "" {
				return ip
			}
		}
	}

	// Fall back to the configured server_url hostname.
	if b.Cfg != nil && b.Cfg.ServerURL != "" {
		// Extract hostname from server_url.
		if idx := strings.Index(b.Cfg.ServerURL, "://"); idx >= 0 {
			rest := b.Cfg.ServerURL[idx+3:]
			if slash := strings.Index(rest, "/"); slash >= 0 {
				rest = rest[:slash]
			}
			if colon := strings.Index(rest, ":"); colon >= 0 {
				rest = rest[:colon]
			}
			if rest != "" {
				return rest
			}
		}
	}

	return ""
}
