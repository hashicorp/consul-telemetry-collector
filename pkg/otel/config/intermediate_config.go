package config

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/service"
)

// IntermediateConfig holds a default list of the ordered components for a pipeline
// It is considered intermediate because in of itself it is never used. Instead
// we build it to have all of the ids of all our generated receivers, exporters,
// and processors. Then when we build a pipeline we clone it and filter out the
// values we decide a pipeline no longer needs.
type IntermediateConfig struct {
	Receivers, Exporters, Processors, Extensions []component.ID
}

// NewIntermediateConfig generates an intermediate configuration used in building up pipelines
func NewIntermediateConfig() *IntermediateConfig {
	return &IntermediateConfig{}
}

// Clone does a deep clone of the intermediate config so we can manipulate it for specific pipelines
// which may want to filter this out
func (i *IntermediateConfig) Clone() *IntermediateConfig {
	newReceivers := make([]component.ID, len(i.Receivers))
	copy(newReceivers, i.Receivers)
	newExporters := make([]component.ID, len(i.Exporters))
	copy(newExporters, i.Exporters)
	newProcessors := make([]component.ID, len(i.Processors))
	copy(newProcessors, i.Processors)
	newExtensions := make([]component.ID, len(i.Extensions))
	copy(newExtensions, i.Extensions)

	return &IntermediateConfig{
		Receivers:  newReceivers,
		Exporters:  newExporters,
		Processors: newProcessors,
		Extensions: newExtensions,
	}
}

// FilterReceiver filters out specified component ids if found
func (i *IntermediateConfig) FilterReceiver(id component.ID) *IntermediateConfig {
	i.Receivers = deleteItemFromSlice(id, i.Receivers)
	return i
}

// FilterExporter filters out specified component ids if found
func (i *IntermediateConfig) FilterExporter(id component.ID) *IntermediateConfig {
	i.Exporters = deleteItemFromSlice(id, i.Exporters)
	return i
}

// FilterProcessor filters out specified component ids if found
func (i *IntermediateConfig) FilterProcessor(id component.ID) *IntermediateConfig {
	i.Processors = deleteItemFromSlice(id, i.Processors)
	return i
}

func deleteItemFromSlice(id component.ID, list []component.ID) []component.ID {
	for idx, v := range list {
		if v == id {
			list = append(list[:idx], list[idx+1:]...)
		}
	}
	return list
}

// ToPipelineConfig translates an intermediate set of ids to a pipeline configuration
// note that pipeline configs do not contain an extension but the service config will
// use the extensions stored in the intermediate config
func (i *IntermediateConfig) ToPipelineConfig() *service.PipelineConfig {
	return &service.PipelineConfig{
		Processors: i.Processors,
		Receivers:  i.Receivers,
		Exporters:  i.Exporters,
	}
}
