package confhelper

import (
	"go.opentelemetry.io/collector/component"

	"github.com/hashicorp/consul-telemetry-collector/pkg/confresolver"
)

func OTLPReceiver(c *confresolver.Config, pipelineIDer confresolver.PipelineIDer,
	pipelines ...confresolver.PipelineIDer) confresolver.ComponentConfig {
	receiver := c.NewReceiver(component.NewID("otlp"), pipelineIDer, pipelines...)
	protocols := receiver.SetMap("protocols")
	protocols.SetMap("http")
	return receiver
}
