// Package config manages helpers to generate opentelemetry-collector configuration.
//
// Currently the Config struct is the overarching configuration that is generated and marshalled to a confmap
// to be used by the providers in the otel/collector package.
//
//	type Config struct {
//	    Receivers  telemetryComponents `mapstructure:"receivers"`
//	    Exporters  telemetryComponents `mapstructure:"exporters"`
//	    Processors telemetryComponents `mapstructure:"processors"`
//	    Connectors telemetryComponents `mapstructure:"connectors"`
//	    Extensions telemetryComponents `mapstructure:"extensions"`
//	    Service    service.Config      `mapstructure:"service"`
//	 }
//
// First off we setup our extensions. Then we define our pipelines and from those definitions we
// can build pipeline configuration.  The service config below shows how the Pipelines are defined in map of id -> pipeline. Each pipeline
// is just a list of receivers, exporters and processors that are dynamically built from that list.
//
//	   Ex: "hcpPipeline"->pipelineConfig of receivers, exporters, processors.
//		/*
//		// Service Configuration
//		type Config struct {
//			Telemetry telemetry.Config `mapstructure:"telemetry"`
//			Extensions []component.ID `mapstructure:"extensions"`
//			Pipelines map[component.ID]*PipelineConfig `mapstructure:"pipelines"`
//		}
//
//		type PipelineConfig struct {
//			Receivers  []component.ID `mapstructure:"receivers"`
//			Processors []component.ID `mapstructure:"processors"`
//			Exporters  []component.ID `mapstructure:"exporters"`
//		}
//		*/
//
// All of this is marshalled to a configuration that the otel collector sdk will run. The goal of this package is to help
// build a configuration that the marshaller will run with our defaults. We setup specific IDs for each pipeline
//
//	/*
//	     receivers:
//			otlp:
//				protocols:
//					http: {}
//
//		processors:
//			memory_limiter:
//				check_interval: 1s
//				limit_percentage: 50
//				spike_limit_percentage: 30
//			batch:
//
//			exporters:
//				logging: {}
//
//			service:
//				pipelines:
//					traces:
//						receivers: [otlp]
//						exporters: [logging]
//	*/
package config
