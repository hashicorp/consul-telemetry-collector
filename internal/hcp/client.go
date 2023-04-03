package hcp

import (
	"errors"
	"fmt"

	"github.com/go-openapi/runtime"
	"github.com/hashicorp/hcp-sdk-go/clients/cloud-global-network-manager-service/preview/2022-02-15/client/global_network_manager_service"
	hcpconfig "github.com/hashicorp/hcp-sdk-go/config"
	"github.com/hashicorp/hcp-sdk-go/httpclient"
	"github.com/hashicorp/hcp-sdk-go/profile"
	"github.com/hashicorp/hcp-sdk-go/resource"
)

// Params is structure used to hold parameters to generate a new client
type Params struct {
	ClientID, ClientSecret, ResourceURL string
}

// telemetryConfig is an internal structure use to store values from the ccm result
// from the api. We use a temporary structure to be able to vary responses between
// versions of the api so we don't have to handle multiple payload versions.
type telemetryConfig struct {
	labels      map[string]string
	endpoint    string
	includeList []string
}

// Client provides a TelemetryClient that lazily retrieves configuration from HCP.
// TelemtryConfiguration can be loaded on-demand using the ReloadConfig() function
type Client struct {
	metricCfg     *telemetryConfig
	hcpResource   *resource.Resource
	clientService agentTelemetryConfigClient
}

var _ TelemetryClient = (*Client)(nil)

// agentTelemetryConfigClient is the interface we expect from the client we
// create. If additional endpoints are needed this interface should expand to
// handle those additional endpoints. Unfortunately the hcp-sdk does not generate
// a mocked client so we use this to build our own mocks as necessary
type agentTelemetryConfigClient interface {
	AgentTelemetryConfig(
		params *global_network_manager_service.AgentTelemetryConfigParams,
		authInfo runtime.ClientAuthInfoWriter,
		opts ...global_network_manager_service.ClientOption,
	) (*global_network_manager_service.AgentTelemetryConfigOK, error)
}

const sourceChannel = "consul-telemetry"

// New creates a new telemetry client for the provided resource using the credentials.
func New(p *Params) (*Client, error) {
	r, err := parseResource(p.ResourceURL)
	if err != nil {
		return nil, err
	}
	hcpConfig, err := parseConfig(p, r)
	if err != nil {
		return nil, err
	}
	runtime, err := httpclient.New(httpclient.Config{
		HCPConfig:     hcpConfig,
		SourceChannel: sourceChannel,
	})
	if err != nil {
		return nil, err
	}
	ccMClient, err := global_network_manager_service.New(runtime, nil), nil
	if err != nil {
		return nil, err
	}
	return newClient(p, ccMClient)
}

// newClient is an internal implementation that takes a clientFn to do deped
func newClient(p *Params, gnmClient agentTelemetryConfigClient) (*Client, error) {
	r, err := parseResource(p.ResourceURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse resource_url %w", err)
	}

	return &Client{
		hcpResource:   r,
		clientService: gnmClient,
	}, nil
}

func parseResource(res string) (*resource.Resource, error) {
	r, err := resource.FromString(res)
	if err != nil {
		return nil, fmt.Errorf("failed to parse resource_url %w", err)
	}

	return &r, nil
}

func parseConfig(p *Params, r *resource.Resource) (hcpconfig.HCPConfig, error) {
	if p.ClientID == "" || p.ClientSecret == "" {
		return nil, errors.New("client credentials are empty")
	}

	hcpconfig, err := hcpconfig.NewHCPConfig(
		hcpconfig.FromEnv(),
		hcpconfig.WithClientCredentials(p.ClientID, p.ClientSecret),
		hcpconfig.WithProfile(&profile.UserProfile{
			OrganizationID: r.Organization,
			ProjectID:      r.Project,
		}),
	)
	if err != nil {
		return nil, errors.New("failed to build hcp config")
	}
	return hcpconfig, nil
}

// ReloadConfig will retrieve the telemetry configuration from HCP using the initially configured runtime.
func (c *Client) ReloadConfig() error {
	params := global_network_manager_service.NewAgentTelemetryConfigParams()
	params.SetClusterID(c.hcpResource.ID)
	result, err := c.clientService.AgentTelemetryConfig(params, nil)
	if err != nil {
		return err
	}
	endpoint := result.Payload.TelemetryConfig.Endpoint
	if result.Payload.TelemetryConfig.Metrics.Endpoint != "" {
		endpoint = result.Payload.TelemetryConfig.Metrics.Endpoint
	}
	metricCfg := telemetryConfig{
		labels:      result.Payload.TelemetryConfig.Labels,
		endpoint:    endpoint,
		includeList: result.Payload.TelemetryConfig.Metrics.IncludeList,
	}

	c.metricCfg = &metricCfg
	return nil
}

// MetricsEndpoint returns the metrics endpoint from the TelemetryConfig
func (c *Client) MetricsEndpoint() (string, error) {
	if c.metricCfg == nil {
		if err := c.ReloadConfig(); err != nil {
			return "", err
		}
	}
	return c.metricCfg.endpoint, nil
}

// MetricFilters returns the metric inclusion filters from the TelemetryConfig
func (c *Client) MetricFilters() ([]string, error) {
	if c.metricCfg == nil {
		if err := c.ReloadConfig(); err != nil {
			return nil, err
		}
	}
	return c.metricCfg.includeList, nil
}
