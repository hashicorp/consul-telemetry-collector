package hcp

type MockClient struct {
	MetricsEndpoint_ string
	MetricFilters_   []string
}

var _ TelemetryClient = (*MockClient)(nil)

func (m *MockClient) MetricsEndpoint() (string, error) {
	return m.MetricsEndpoint_, nil
}

func (m *MockClient) MetricFilters() ([]string, error) {
	return m.MetricFilters_, nil
}
