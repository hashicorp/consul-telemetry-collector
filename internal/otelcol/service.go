package otelcol

import (
	"context"
	"fmt"

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
	factories, err := components()
	if err != nil {
		return nil, err
	}

	cfgProvider, err := Provider()
	if err != nil {
		return nil, err
	}

	staticCfg, err := cfgProvider.Get(ctx, factories)
	if err != nil {
		fmt.Println(err)
	}

	cfg := service.Config{
		Telemetry: telemetry.Config{
			Logs: telemetry.LogsConfig{
				Level:       zapcore.DebugLevel,
				Development: true,
				Encoding:    "console",
				OutputPaths: []string{"stderr"},
			},
			Metrics: staticCfg.Service.Telemetry.Metrics,
			Traces: telemetry.TracesConfig{
				Propagators: []string{},
			},
			Resource: map[string]*string{
				"": nil,
			},
		},
		Pipelines: staticCfg.Service.Pipelines,
	}

	set := service.Settings{
		BuildInfo: component.BuildInfo{
			Command:     "",
			Description: "",
			Version:     "",
		},
		Receivers:         receiver.NewBuilder(staticCfg.Receivers, factories.Receivers),
		Processors:        &processor.Builder{},
		Exporters:         exporter.NewBuilder(staticCfg.Exporters, factories.Exporters),
		Connectors:        &connector.Builder{},
		Extensions:        &extension.Builder{},
		AsyncErrorChannel: make(chan error),
		LoggingOptions:    []zap.Option{},
	}

	return service.New(ctx, set, cfg)
}
