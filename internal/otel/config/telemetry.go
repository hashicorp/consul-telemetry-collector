package config

import (
	"go.opentelemetry.io/collector/config/configtelemetry"
	"go.opentelemetry.io/collector/service/telemetry"
	"go.uber.org/zap/zapcore"
)

const ( // supported trace propagators
	traceContextPropagator = "tracecontext"
	b3Propagator           = "b3"
)

// Telemetry returns our basic telemetry configuration.
func Telemetry() telemetry.Config {
	return telemetry.Config{
		Logs: telemetry.LogsConfig{
			Level:       zapcore.InfoLevel,
			Encoding:    "console",
			OutputPaths: []string{"stderr"},
		},
		Metrics: telemetry.MetricsConfig{
			Address: "localhost:9090",
			Level:   configtelemetry.LevelDetailed,
		},
		Traces: telemetry.TracesConfig{
			Propagators: []string{traceContextPropagator, b3Propagator},
		},
	}
}
