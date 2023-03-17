package hcp

type MockClient struct {
	HCPMetricsEndpoint string
}

var _ TelemetryClient = (*MockClient)(nil)

func (m *MockClient) MetricsEndpoint() string {
	return m.HCPMetricsEndpoint
}
