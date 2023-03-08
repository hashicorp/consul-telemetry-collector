// Package confresolver manages helpers to generate opentelemetry-collector configuration.
// The helpers force all component configuration to get added to a pipeline.
// After initial creation a pipeline is created and returns a pipeline reference.
// The pipeline reference is a custom type that allows us to retrieve the pipeline when creating new exporters and
// receivers. This ensures that any telemetryComponents that we're configuring also get added to a pipeline.
//
// The helpers ensure that all components when created are included in a pipeline and telemetryComponents configuration
// For example, the following code would create the equivalent of this yaml configuration.
// In this example the `otlp` receiver and `logging` exporter components are included in the appropriate component
// configuration sections and also in the Traces `pipeline`.
// This ensures that the configured components are always used in the pipeline.
//
//	pipeline := c.NewPipeline(component.DataTypeTraces)
//	receiver := c.NewReceiver(pipeline, component.NewID("otlp"))
//	receiver.Map("protocols").Map("http")
//	c.NewExporter(pipeline, component.NewID("logging"))
//
//	/*
//			receivers:
//				otlp:
//					protocols:
//			  			http: {}
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
package confresolver
