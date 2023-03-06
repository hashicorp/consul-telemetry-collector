package otelcol

import (
	"context"

	"go.opentelemetry.io/collector/service"
	"go.opentelemetry.io/collector/service/telemetry"
)

// New will create a new open-telemetry collector service and configuration based on the provided values
func New(ctx context.Context) (*service.Service, error) {
	cfg := service.Config{
		Telemetry: telemetry.Config{
			Logs: telemetry.LogsConfig{
				Encoding:    "console",
				OutputPaths: []string{"stderr"},
			},
			Metrics: telemetry.MetricsConfig{},
			Traces:  telemetry.TracesConfig{},
		},
	}

	settings := service.Settings{}

	return service.New(ctx, settings, cfg)
}
