package otelcol

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/connector"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/loggingexporter"
	"go.opentelemetry.io/collector/extension"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
	"go.opentelemetry.io/collector/service"
	"go.opentelemetry.io/collector/service/telemetry"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Receivers() receiver.Factory {
	return otlpreceiver.NewFactory()
}

func Exporters() exporter.Factory {
	return loggingexporter.NewFactory()
}

func NewService(ctx context.Context, _ service.Settings, _ service.Config) (*service.Service, error) {
	traces := component.NewID(component.DataTypeTraces)

	r := Receivers()
	e := Exporters()

	factories, err := components()
	if err != nil {
		return nil, err
	}

	cfg := service.Config{
		Telemetry: telemetry.Config{
			Logs: telemetry.LogsConfig{
				Level:       zapcore.DebugLevel,
				Development: true,
				Encoding:    "console",
				OutputPaths: []string{"stderr"},
			},
			Metrics: telemetry.MetricsConfig{
				Level:   0,
				Address: "",
			},
			Traces: telemetry.TracesConfig{
				Propagators: []string{},
			},
			Resource: map[string]*string{
				"": nil,
			},
		},
		Pipelines: map[component.ID]*service.PipelineConfig{
			traces: {
				Receivers:  []component.ID{component.NewID(r.Type())},
				Processors: []component.ID{},
				Exporters:  []component.ID{component.NewID(e.Type())},
			},
		},
	}

	set := service.Settings{
		BuildInfo: component.BuildInfo{
			Command:     "",
			Description: "",
			Version:     "",
		},
		Receivers: receiver.NewBuilder(map[component.ID]component.Config{
			component.NewID(r.Type()): r.CreateDefaultConfig(),
		}, factories.Receivers),
		Processors:        &processor.Builder{},
		Exporters:         exporter.NewBuilder(map[component.ID]component.Config{component.NewID(e.Type()): e.CreateDefaultConfig()}, factories.Exporters),
		Connectors:        &connector.Builder{},
		Extensions:        &extension.Builder{},
		AsyncErrorChannel: make(chan error),
		LoggingOptions:    []zap.Option{},
	}

	return service.New(ctx, set, cfg)
}
