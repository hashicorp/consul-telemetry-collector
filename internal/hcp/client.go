package hcp

import (
	"errors"
	"fmt"

	"github.com/go-openapi/runtime/client"

	"github.com/hashicorp/hcp-sdk-go/clients/cloud-global-network-manager-service/preview/2022-02-15/client/global_network_manager_service"
	hcpconfig "github.com/hashicorp/hcp-sdk-go/config"
	"github.com/hashicorp/hcp-sdk-go/httpclient"
	"github.com/hashicorp/hcp-sdk-go/profile"
	"github.com/hashicorp/hcp-sdk-go/resource"
)

const sourceChannel = "consul-telemetry"

// Client provides a TelemetryClient that lazily retrieves configuration from HCP.
// TelemtryConfiguration can be loaded on-demand using the ReloadConfig() function
type Client struct {
	runtime     *client.Runtime
	metricCfg   *telemetryConfig
	hcpResource resource.Resource
	ccmClient   ClientService
}

type telemetryConfig struct {
	labels      map[string]string
	endpoint    string
	includeList []string
}

var _ TelemetryClient = (*Client)(nil)

// New creates a new telemetry client for the provided resource using the credentials.
func New(clientID, clientSecret, resourceURL string) (*Client, error) {
	r, err := resource.FromString(resourceURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse resource_url %w", err)
	}

	if clientID == "" || clientSecret == "" {
		return nil, errors.New("client credentials are empty")
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

type clientOpt func(*Client)

// NewWithDeps creates a new client while also providing client modification via client
func NewWithDeps(clientID, clientSecret, resourceURL string, opts ...clientOpt) (*Client, error) {
	client, err := New(clientID, clientSecret, resourceURL)
	if err != nil {
		return nil, err
	}
	for _, opt := range opts {
		opt(client)
	}

	return client, nil
}

// WithClientService sets the provided ClientService on the Client
//
//revive:disable:unexported-return
func WithClientService(cs ClientService) clientOpt {
	return clientOpt(func(c *Client) {
		c.ccmClient = cs
	})
}

//revive:enable:unexported-return

// ReloadConfig will retrieve the telemetry configuration from HCP using the initially configured runtime.
func (c *Client) ReloadConfig() error {
	metricsCfg, err := getTelemetryConfig(c.ccmClient, c.hcpResource.ID)
	if err != nil {
		return err
	}
	c.metricCfg = &metricsCfg
	return nil
}

func getTelemetryConfig(gnm ClientService, clusterID string) (telemetryConfig, error) {
	params := global_network_manager_service.NewAgentTelemetryConfigParams()
	params.SetClusterID(clusterID)
	result, err := gnm.AgentTelemetryConfig(params, nil)
	if err != nil {
		return telemetryConfig{}, err
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

	return metricCfg, nil
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
