package confresolver

import "go.opentelemetry.io/collector/component"

// NewReceiver creates a new receiver configuration and adds it to the referenced pipeline
func (c *Config) NewReceiver(p PipelineIDer, id component.ID) ComponentConfig {
	var ccfg componentConfig

	if c.Receivers == nil {
		c.Receivers = make(components)
	}

	pipeline := c.Service.Pipelines[p.id()]
	pipeline.Receivers, ccfg = addComponent(c.Receivers, id, pipeline.Receivers)

	return ccfg
}
