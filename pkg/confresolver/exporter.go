package confresolver

import "go.opentelemetry.io/collector/component"

// NewExporter adds a new exporter component to the overall configuration and
// adds the component to the specified pipeline. It returns the component's
// configuration for further configuring. The PipelineIDer parameter is the reference
// // returned from NewPipeline
func (c *Config) NewExporter(p PipelineIDer, id component.ID) ComponentConfig {
	var ccfg componentConfig

	if c.Exporters == nil {
		c.Exporters = make(telemetryComponents)
	}

	pipeline := c.Service.Pipelines[p.id()]
	pipeline.Exporters, ccfg = addComponent(c.Exporters, id, pipeline.Exporters)

	return ccfg
}
