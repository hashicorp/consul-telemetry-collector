package confresolver

import "go.opentelemetry.io/collector/component"

// NewExporter creates a new component configuration and adds it to the specified pipeline
func (c *Config) NewExporter(p PipelineIDer, id component.ID) ComponentConfig {
	var ccfg componentConfig

	if c.Exporters == nil {
		c.Exporters = make(components)
	}

	pipeline := c.Service.Pipelines[p.id()]
	pipeline.Exporters, ccfg = addComponent(c.Exporters, id, pipeline.Exporters)

	return ccfg
}
