package confhelper

import (
	"go.opentelemetry.io/collector/component"

	"github.com/hashicorp/consul-telemetry-collector/pkg/confresolver"
)

const (
	memoryLimiterID = "memory_limiter"
)

// MemoryLimiter creates a process that ensures that the memory utilization of the open-telemetry-collector doesn't
// go above 50% of the total available memory with a 30% burst
func MemoryLimiter(c *confresolver.Config, pipelineIDer confresolver.PipelineIDer,
	pipelines ...confresolver.PipelineIDer) confresolver.ComponentConfig {

	const (
		checkInterval        = "check_interval"
		limitPercentage      = "limit_percentage"
		spikeLimitPercentage = "spike_limit_percentage"
	)

	limiter := c.NewProcessor(component.NewID(memoryLimiterID), pipelineIDer, pipelines...)
	limiter.Set(checkInterval, "1s")
	limiter.Set(limitPercentage, "50")
	limiter.Set(spikeLimitPercentage, "30")
	return limiter
}
