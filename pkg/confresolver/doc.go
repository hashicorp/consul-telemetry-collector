// Package confresolver manages helpers to generate opentelemetry-collector configuration.
// The helpers force all component configuration to get added to a pipeline.
// After initial creation a pipeline is created and returns a pipeline reference.
// The pipeline reference is a custom type that allows us to retrieve the pipeline when creating new exporters and
// receivers. This ensures that any components that we're configuring also get added to a pipeline.
package confresolver
