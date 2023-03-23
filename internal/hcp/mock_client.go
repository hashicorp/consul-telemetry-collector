package hcp

// MockClient fulfills the TelemetryClient interface and returns static values. Used for testing
type MockClient struct {
	MockMetricsEndpoint string
	MockMetricFilters   []string
}

var _ TelemetryClient = (*MockClient)(nil)

// MetricsEndpoint returns the provided metrics endpoint. Will never error.
func (m *MockClient) MetricsEndpoint() (string, error) {
	return m.MockMetricsEndpoint, nil
}

// MetricFilters returns the provided metric inclusion filters. Will never error.
func (m *MockClient) MetricFilters() ([]string, error) {
	return m.MockMetricFilters, nil
}
