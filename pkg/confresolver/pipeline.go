package confresolver

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/service"
)

// PipelineIDer acts as a reference to a named or unnamed pipeline
type PipelineIDer interface {
	id() component.ID
}

type pipelineRef component.ID

func (p pipelineRef) id() component.ID {
	return component.ID(p)
}

// NewPipeline creates a new, unnamed pipeline in the configuration
func (c *Config) NewPipeline(pipeline component.DataType) PipelineIDer {
	if c.Service.Pipelines == nil {
		c.Service.Pipelines = make(map[component.ID]*service.PipelineConfig)
	}

	id := component.NewID(pipeline)

	c.Service.Pipelines[id] = &service.PipelineConfig{}
	return pipelineRef(id)
}

// NewPipelineWithName creates a new pipeline with a specified name.
func (c *Config) NewPipelineWithName(pipeline component.DataType, name string) PipelineIDer {
	if c.Service.Pipelines == nil {
		c.Service.Pipelines = make(map[component.ID]*service.PipelineConfig)
	}

	id := component.NewIDWithName(pipeline, name)

	c.Service.Pipelines[id] = &service.PipelineConfig{}
	return pipelineRef(id)
}

func (c *Config) PushExporterOnPipeline(p PipelineIDer, id ...component.ID) {
	c.Service.Pipelines[p.id()].Exporters = append(c.Service.Pipelines[p.id()].Exporters, id...)
}

func (c *Config) PushProcessorOnPipeline(p PipelineIDer, id ...component.ID) {
	c.Service.Pipelines[p.id()].Processors = append(c.Service.Pipelines[p.id()].Processors, id...)
}

func (c *Config) PushReceiverOnPipeline(p PipelineIDer, id ...component.ID) {
	c.Service.Pipelines[p.id()].Receivers = append(c.Service.Pipelines[p.id()].Receivers, id...)
}
