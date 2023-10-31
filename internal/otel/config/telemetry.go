// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"fmt"

	"go.opentelemetry.io/collector/config/configtelemetry"
	"go.opentelemetry.io/collector/service/telemetry"
	"go.uber.org/zap/zapcore"
)

const ( // supported trace propagators
	traceContextPropagator = "tracecontext"
	b3Propagator           = "b3"
)

// Telemetry returns our basic telemetry configuration.
func Telemetry(metricsPort int) telemetry.Config {
	return telemetry.Config{
		Logs: telemetry.LogsConfig{
			Level:            zapcore.InfoLevel,
			Encoding:         "console",
			OutputPaths:      []string{"stderr"},
			ErrorOutputPaths: []string{"stderr"},
		},
		Metrics: telemetry.MetricsConfig{
			Address: fmt.Sprintf("localhost:%d", metricsPort),
			Level:   configtelemetry.LevelDetailed,
			Readers: []telemetry.MetricReader{},
		},
		Traces: telemetry.TracesConfig{
			Propagators: []string{traceContextPropagator, b3Propagator},
			Processors:  []telemetry.SpanProcessor{},
		},
	}
}
