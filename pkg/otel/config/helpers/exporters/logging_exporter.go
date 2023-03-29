package exporters

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configtelemetry"
	"go.opentelemetry.io/collector/exporter/loggingexporter"
)

// LoggingConfig is temporary until we can use the config under this
// module:
//
//	"go.opentelemetry.io/collector/exporter/loggingexporter"
//
// A bug exists that I have patched here but it's pending on CLA acceptance
//   - PR: https://github.com/open-telemetry/opentelemetry-collector/pull/7447
type LoggingConfig struct {
	// Verbosity defines the logging exporter verbosity.
	Verbosity configtelemetry.Level `mapstructure:"verbosity"`

	// SamplingInitial defines how many samples are initially logged during each second.
	SamplingInitial int `mapstructure:"sampling_initial"`

	// SamplingThereafter defines the sampling rate after the initial samples are logged.
	SamplingThereafter int `mapstructure:"sampling_thereafter"`
}

// loggingExporterName is the component.ID value used by the logging exporter
const loggingExporterName = "logging"

// LoggingExporterID is the component.ID value used by the logging exporter
var LoggingExporterID component.ID = component.NewID(loggingExporterName)

// LogExporterCfg generates the configuration for a logging exporter
func LogExporterCfg() (component.ID, *LoggingConfig) {
	defaults := loggingexporter.NewFactory().CreateDefaultConfig().(*loggingexporter.Config)

	return LoggingExporterID, &LoggingConfig{
		Verbosity:          defaults.Verbosity,
		SamplingInitial:    defaults.SamplingInitial,
		SamplingThereafter: defaults.SamplingThereafter,
	}
}
