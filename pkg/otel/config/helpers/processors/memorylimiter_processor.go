package processors

import (
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/processor/memorylimiterprocessor"
)

const memoryLimiterName = "memory_limiter"

// MemoryLimiterID is the component id of the memory limiter
var MemoryLimiterID component.ID = component.NewID(memoryLimiterName)

// MemoryLimiterCfg  generates the config for a memory limiter processor
func MemoryLimiterCfg() *memorylimiterprocessor.Config {
	return &memorylimiterprocessor.Config{
		CheckInterval:         time.Second,
		MemoryLimitPercentage: 50,
		MemorySpikePercentage: 30,
	}
}
