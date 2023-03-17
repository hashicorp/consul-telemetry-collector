package hcp

type MockClient struct {
	MetricsEndpoint_ string
	MetricFilters_   []string
}

var _ TelemetryClient = (*MockClient)(nil)

func (m *MockClient) MetricsEndpoint() string {
	return m.MetricsEndpoint_
}

func (m *MockClient) MetricFilters() []string {
	return m.MetricFilters_
}
