package confresolver

import "go.opentelemetry.io/collector/component"

// NewReceiver adds a new receiver component to the overall configuration and
// adds the component to the specified pipeline. It returns the component's
// configuration for further configuring. The PipelineIDer parameter is the reference
// returned from NewPipeline
func (c *Config) NewReceiver(id component.ID, pipelineIDer PipelineIDer, pipelineIDs ...PipelineIDer) ComponentConfig {
	var ccfg componentConfig

	if c.Receivers == nil {
		c.Receivers = make(telemetryComponents)
	}

	pipeline := c.Service.Pipelines[pipelineIDer.id()]
	pipeline.Receivers, ccfg = addComponent(c.Receivers, id, pipeline.Receivers)

	for _, p := range pipelineIDs {
		pipeline := c.Service.Pipelines[p.id()]
		pipeline.Receivers = append(pipeline.Receivers, id)
	}

	return ccfg
}
