package confresolver

import "go.opentelemetry.io/collector/component"

// NewExporter adds a new exporter component to the overall configuration and
// adds the component to the specified pipeline. It returns the component's
// configuration for further configuring. The PipelineIDer parameter is the reference
// // returned from NewPipeline
func (c *Config) NewExporter(id component.ID, pipelineIDer PipelineIDer, pipelineIDs ...PipelineIDer) ComponentConfig {
	var ccfg componentConfig

	if c.Exporters == nil {
		c.Exporters = make(telemetryComponents)
	}

	pipeline := c.Service.Pipelines[pipelineIDer.id()]
	pipeline.Exporters, ccfg = addComponent(c.Exporters, id, pipeline.Exporters)

	for _, p := range pipelineIDs {
		pipeline := c.Service.Pipelines[p.id()]
		pipeline.Exporters = append(pipeline.Exporters, id)
	}

	return ccfg
}
