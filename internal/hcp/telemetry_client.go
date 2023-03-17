package hcp

import (
	"github.com/go-openapi/runtime"
	"github.com/hashicorp/hcp-sdk-go/clients/cloud-global-network-manager-service/preview/2022-02-15/client/global_network_manager_service"
)

// TelemetryClient is a high level client for the AgentTelemetryConfig.
// It abstracts the interaction with HCP to retrieve the AgentTelemetryConfig.
type TelemetryClient interface {
	MetricsEndpoint() (string, error)
	MetricFilters() ([]string, error)
}

// ClientService is a paired down interface for the global-network-manager-service that retrieves the
// AgentTelemetryConfig. Allows mocking
type ClientService interface {
	AgentTelemetryConfig(params *global_network_manager_service.AgentTelemetryConfigParams, authInfo runtime.ClientAuthInfoWriter,
		opts ...global_network_manager_service.ClientOption) (*global_network_manager_service.AgentTelemetryConfigOK,
		error)
}
