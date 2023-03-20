package confhelper

import (
	"go.opentelemetry.io/collector/component"

	"github.com/hashicorp/consul-telemetry-collector/pkg/confresolver"
)

// MemoryLimiter creates a process that ensures that the memory utilization of the open-telemetry-collector doesn't
// go above 50% of the total available memory with a 30% burst
func MemoryLimiter(c *confresolver.Config, pipelineIDer confresolver.PipelineIDer,
	pipelines ...confresolver.PipelineIDer) confresolver.ComponentConfig {

	limiter := c.NewProcessor(component.NewID("memory_limiter"), pipelineIDer, pipelines...)
	limiter.Set("check_interval", "1s")
	limiter.Set("limit_percentage", "50")
	limiter.Set("spike_limit_percentage", "30")
	return limiter
}
