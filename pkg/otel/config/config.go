package config

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/service"
)

// Config is a helper type to create a new opentelemetry server configuration.
// It implements a map[string]interface{} representation of the opentelemetry-collector configuration.
// More information can be found here: https://opentelemetry.io/docs/collector/configuration/
type Config struct {
	Receivers  telemetryComponents `mapstructure:"receivers"`
	Exporters  telemetryComponents `mapstructure:"exporters"`
	Processors telemetryComponents `mapstructure:"processors"`
	Connectors telemetryComponents `mapstructure:"connectors"`
	Extensions telemetryComponents `mapstructure:"extensions"`
	Service    service.Config      `mapstructure:"service"`
}

// NewConfig creates a new config object with all types initialized
func NewConfig() *Config {
	svcConfig := service.Config{}
	svcConfig.Pipelines = make(map[component.ID]*service.PipelineConfig)
	return &Config{
		Receivers:  make(telemetryComponents),
		Exporters:  make(telemetryComponents),
		Processors: make(telemetryComponents),
		Connectors: make(telemetryComponents),
		Extensions: make(telemetryComponents),
		Service:    svcConfig,
	}
}

type telemetryComponents map[component.ID]any
