package hcp

import (
	"github.com/go-openapi/runtime"
	"github.com/hashicorp/hcp-sdk-go/clients/cloud-global-network-manager-service/preview/2022-02-15/client/global_network_manager_service"
)

type TelemetryClient interface {
	MetricsEndpoint() (string, error)
	MetricFilters() ([]string, error)
}

type ClientService interface {
	AgentTelemetryConfig(params *global_network_manager_service.AgentTelemetryConfigParams, authInfo runtime.ClientAuthInfoWriter,
		opts ...global_network_manager_service.ClientOption) (*global_network_manager_service.AgentTelemetryConfigOK,
		error)
}
