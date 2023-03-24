package confresolver

import "go.opentelemetry.io/collector/component"

// NewProcessor adds a new processor component to the overall configuration and
// adds the component to the specified pipeline. It returns the component's
// configuration for further configuring. The PipelineIDer parameter is the reference
// returned from NewPipeline
func (c *Config) NewProcessor(id component.ID, pipelineIDer PipelineIDer, pipelineIDs ...PipelineIDer) ComponentConfig {
	var ccfg componentConfig

	if c.Processors == nil {
		c.Processors = make(telemetryComponents)
	}

	pipeline := c.Service.Pipelines[pipelineIDer.id()]
	pipeline.Processors, ccfg = addComponent(c.Processors, id, pipeline.Processors)

	for _, p := range pipelineIDs {
		pipeline := c.Service.Pipelines[p.id()]
		pipeline.Processors = append(pipeline.Processors, id)
	}

	return ccfg
}
