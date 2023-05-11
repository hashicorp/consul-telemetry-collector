package hcp

import (
	"github.com/go-openapi/runtime"

	"github.com/hashicorp/hcp-sdk-go/clients/cloud-consul-telemetry-gateway/preview/2023-04-14/client/consul_telemetry_service"
)

// MockClientService fulfills the TelemetryClient interface and returns static values. Used for testing.
type MockClientService struct {
	MockResponse *consul_telemetry_service.AgentTelemetryConfigOK
	Err          error
	params       *consul_telemetry_service.AgentTelemetryConfigParams
	opts         []consul_telemetry_service.ClientOption
}

var _ ClientService = (*MockClientService)(nil)

// AgentTelemetryConfig returns mocked responses.
func (m *MockClientService) AgentTelemetryConfig(params *consul_telemetry_service.AgentTelemetryConfigParams,
	_ runtime.ClientAuthInfoWriter,
	opts ...consul_telemetry_service.ClientOption) (*consul_telemetry_service.AgentTelemetryConfigOK,
	error) {
	m.params = params
	m.opts = opts
	return m.MockResponse, m.Err
}
