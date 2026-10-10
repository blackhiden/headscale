package apiv1

import (
	"context"
	"net/http"
	"slices"

	"github.com/danielgtaylor/huma/v2"
	"tailscale.com/tailcfg"
)

// DERPNode mirrors a single DERP relay node in the DERP map.
type DERPNode struct {
	Name     string `json:"name"`
	HostName string `json:"hostName"`
	DERPPort int    `json:"derpPort"`
	STUNPort int    `json:"stunPort"`
	IPv4     string `json:"ipv4"`
	IPv6     string `json:"ipv6"`
}

// DERPRegion mirrors a DERP region (a set of relay nodes).
type DERPRegion struct {
	RegionID   int        `json:"regionId"`
	RegionName string     `json:"regionName"`
	RegionCode string     `json:"regionCode"`
	Nodes      []DERPNode `json:"nodes"`
}

// DERPResponseBody is the response for GET /api/v1/derp.
type DERPResponseBody struct {
	Configured   bool         `json:"configured"`
	TotalRegions int          `json:"totalRegions"`
	Regions      []DERPRegion `json:"regions"`
}

type derpOutput struct {
	Body DERPResponseBody
}

func init() {
	registrations = append(registrations, registerDERP)
}

func registerDERP(api huma.API, b Backend) {
	huma.Register(api, huma.Operation{
		OperationID: "getDerp",
		Method:      http.MethodGet,
		Path:        "/api/v1/derp",
		Summary:     "Get DERP map",
		Description: "Returns the current DERP relay map configuration, including regions and their nodes.",
		Tags:        []string{"DERP"},
		Security:    bearerAuth,
	}, func(ctx context.Context, _ *struct{}) (*derpOutput, error) {
		derpMap := b.State.DERPMap()

		body := DERPResponseBody{
			Configured: derpMap.Valid(),
			Regions:    []DERPRegion{},
		}

		if !derpMap.Valid() {
			return &derpOutput{Body: body}, nil
		}

		regionIDs := make([]tailcfg.DERPRegionID, 0, derpMap.Regions().Len())
		for regionID := range derpMap.Regions().All() {
			regionIDs = append(regionIDs, regionID)
		}

		slices.Sort(regionIDs)

		for _, regionID := range regionIDs {
			region := derpMap.Regions().Get(regionID)
			apiRegion := DERPRegion{
				RegionID:   int(region.RegionID()),
				RegionName: region.RegionName(),
				RegionCode: region.RegionCode(),
				Nodes:      []DERPNode{},
			}

			for _, node := range region.Nodes().All() {
				apiRegion.Nodes = append(apiRegion.Nodes, DERPNode{
					Name:     node.Name(),
					HostName: node.HostName(),
					DERPPort: node.DERPPort(),
					STUNPort: node.STUNPort(),
					IPv4:     node.IPv4(),
					IPv6:     node.IPv6(),
				})
			}

			body.Regions = append(body.Regions, apiRegion)
		}

		body.TotalRegions = len(body.Regions)

		return &derpOutput{Body: body}, nil
	})
}
