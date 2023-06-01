package hcp

// MockClient fulfills the TelemetryClient interface and returns static values. Used for testing.
type MockClient struct {
	MockMetricsEndpoint  string
	MockMetricFilters    []string
	MockMetricAttributes map[string]string
	Err                  error
}

var _ TelemetryClient = (*MockClient)(nil)

// MetricsEndpoint returns the provided metrics endpoint. Will never error.
func (m *MockClient) MetricsEndpoint() (string, error) {
	if m.Err != nil {
		return "", m.Err
	}
	return m.MockMetricsEndpoint, nil
}

// MetricFilters returns the provided metric inclusion filters. Will never error.
func (m *MockClient) MetricFilters() ([]string, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.MockMetricFilters, nil
}

// MetricAttributes returns the provided metric inclusion filters. Will never error.
func (m *MockClient) MetricAttributes() (map[string]string, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.MockMetricAttributes, nil
}
