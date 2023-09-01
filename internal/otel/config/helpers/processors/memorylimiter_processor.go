// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package processors

import (
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/processor/memorylimiterprocessor"
)

const memoryLimiterName = "memory_limiter"

// MemoryLimiterID is the component id of the memory limiter.
var MemoryLimiterID component.ID = component.NewID(memoryLimiterName)

// MemoryLimiterCfg  generates the config for a memory limiter processor.
func MemoryLimiterCfg() *memorylimiterprocessor.Config {
	return &memorylimiterprocessor.Config{
		CheckInterval:         time.Second,
		MemoryLimitPercentage: 80,
		MemorySpikePercentage: 20,
	}
}
