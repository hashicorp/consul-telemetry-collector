package confhelper

import (
	"go.opentelemetry.io/collector/component"

	"github.com/hashicorp/consul-telemetry-collector/pkg/confresolver"
)

const otlpReceiverID = "otlp"

// OTLPReceiver confhelper creates an otlp receiver on the provided pipelines
func OTLPReceiver(c *confresolver.Config, pipelineIDer confresolver.PipelineIDer,
	pipelines ...confresolver.PipelineIDer) confresolver.ComponentConfig {
	receiver := c.NewReceiver(component.NewID(otlpReceiverID), pipelineIDer, pipelines...)
	protocols := receiver.SetMap("protocols")
	protocols.SetMap("http")
	return receiver
}
