package hcp

import (
	"fmt"

	"github.com/go-openapi/runtime/client"
	hcpconfig "github.com/hashicorp/hcp-sdk-go/config"
	"github.com/hashicorp/hcp-sdk-go/httpclient"
	"github.com/hashicorp/hcp-sdk-go/profile"
	"github.com/hashicorp/hcp-sdk-go/resource"
)

const sourceChannel = "consul-telemetry"

type Client struct {
	runtime *client.Runtime
}

var _ TelemetryClient = (*Client)(nil)

func New(clientID, clientSecret, resourceURL string) (*Client, error) {
	r, err := resource.FromString(resourceURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse resource_url %w", err)
	}

	hcpConfig, err := hcpconfig.NewHCPConfig(
		hcpconfig.FromEnv(),
		hcpconfig.WithClientCredentials(clientID, clientSecret),
		hcpconfig.WithProfile(&profile.UserProfile{
			OrganizationID: r.Organization,
			ProjectID:      r.Project,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to configure hcp-sdk client %w", err)
	}

	runtime, err := httpclient.New(httpclient.Config{
		HCPConfig:     hcpConfig,
		SourceChannel: sourceChannel,
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		runtime: runtime,
	}, nil
}

func (c *Client) MetricsEndpoint() string {
	return ""
}

func (c *Client) MetricFilters() []string {
	return []string{}
}
