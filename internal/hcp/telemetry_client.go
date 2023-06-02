package hcp

import (
	"github.com/go-openapi/runtime"

	"github.com/hashicorp/hcp-sdk-go/clients/cloud-consul-telemetry-gateway/preview/2023-04-14/client/consul_telemetry_service"
)

// TelemetryClient is a high level client for the AgentTelemetryConfig.
// It abstracts the interaction with HCP to retrieve the AgentTelemetryConfig.
type TelemetryClient interface {
	MetricsEndpoint() (string, error)
	MetricFilters() ([]string, error)
	MetricAttributes() (map[string]string, error)
}

// ClientService is a paired down interface for the global-network-manager-service that retrieves the
// AgentTelemetryConfig. Allows mocking.
type ClientService interface {
	AgentTelemetryConfig(params *consul_telemetry_service.AgentTelemetryConfigParams, authInfo runtime.ClientAuthInfoWriter,
		opts ...consul_telemetry_service.ClientOption) (*consul_telemetry_service.AgentTelemetryConfigOK,
		error)
}
