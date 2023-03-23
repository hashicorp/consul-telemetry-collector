package hcp

import (
	"github.com/go-openapi/runtime"

	"github.com/hashicorp/hcp-sdk-go/clients/cloud-global-network-manager-service/preview/2022-02-15/client/global_network_manager_service"
)

// MockClientService fulfills the TelemetryClient interface and returns static values. Used for testing
type MockClientService struct {
	MockResponse *global_network_manager_service.AgentTelemetryConfigOK
	Err          error
	params       *global_network_manager_service.AgentTelemetryConfigParams
	opts         []global_network_manager_service.ClientOption
}

var _ ClientService = (*MockClientService)(nil)

// AgentTelemetryConfig returns mocked responses
func (m *MockClientService) AgentTelemetryConfig(params *global_network_manager_service.AgentTelemetryConfigParams, authInfo runtime.ClientAuthInfoWriter,
	opts ...global_network_manager_service.ClientOption) (*global_network_manager_service.AgentTelemetryConfigOK,
	error) {
	m.params = params
	m.opts = opts
	return m.MockResponse, m.Err
}
