package hcp

type TelemetryClient interface {
	MetricsEndpoint() string
}
