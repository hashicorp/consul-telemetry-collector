package hcp

import (
	"fmt"

	"github.com/go-openapi/runtime/client"
	"github.com/hashicorp/hcp-sdk-go/clients/cloud-global-network-manager-service/preview/2022-02-15/client/global_network_manager_service"
	hcpconfig "github.com/hashicorp/hcp-sdk-go/config"
	"github.com/hashicorp/hcp-sdk-go/httpclient"
	"github.com/hashicorp/hcp-sdk-go/profile"
	"github.com/hashicorp/hcp-sdk-go/resource"
)

const sourceChannel = "consul-telemetry"

type Client struct {
	runtime     *client.Runtime
	metricCfg   *telemetryConfig
	hcpResource resource.Resource
}

type telemetryConfig struct {
	labels      map[string]string
	endpoint    string
	includeList []string
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

func (c *Client) LoadTelemetryConfig(gnm ClientService) error {
	metricsCfg, err := GetTelemetryConfig(gnm, c.hcpResource.ID)
	if err != nil {
		return err
	}
	c.metricCfg = &metricsCfg
	return nil
}

func (c *Client) reloadConfig() error {
	gnmClient := global_network_manager_service.New(c.runtime, nil)
	return c.LoadTelemetryConfig(gnmClient)
}

func GetTelemetryConfig(gnm ClientService, clusterID string) (telemetryConfig, error) {
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

func (c *Client) MetricsEndpoint() (string, error) {
	if c.metricCfg == nil {
		if err := c.reloadConfig(); err != nil {
			return "", err
		}
	}
	return c.metricCfg.endpoint, nil
}

func (c *Client) MetricFilters() ([]string, error) {
	if c.metricCfg == nil {
		if err := c.reloadConfig(); err != nil {
			return nil, err
		}
	}
	return c.metricCfg.includeList, nil
}
