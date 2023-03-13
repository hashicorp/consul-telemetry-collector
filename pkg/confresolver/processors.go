package confresolver

import "go.opentelemetry.io/collector/component"

// Processor adds a new processor component to the overall configuration and
// adds the component to the specified pipeline. It returns the component's
// configuration for further configuring. The PipelineIDer parameter is the reference
// returned from NewPipeline
func (c *Config) NewProcessor(p PipelineIDer, id component.ID) ComponentConfig {
	var ccfg componentConfig

	if c.Processors == nil {
		c.Processors = make(telemetryComponents)
	}

	pipeline := c.Service.Pipelines[p.id()]
	pipeline.Processors, ccfg = addComponent(c.Processors, id, pipeline.Processors)

	return ccfg
}
