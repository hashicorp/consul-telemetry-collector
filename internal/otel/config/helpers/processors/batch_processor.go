// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package processors

import (
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/processor/batchprocessor"
)

const batchProcessorName = "batch"

// BatchProcessorID is the component id of the batch processor.
var BatchProcessorID component.ID = component.NewID(batchProcessorName)

// BatchProcessorCfg  generates the config for a batch processor.
func BatchProcessorCfg(timeout time.Duration) *batchprocessor.Config {
	if timeout == 0 {
		timeout = time.Minute
	}
	factory := batchprocessor.NewFactory()
	cfg := factory.CreateDefaultConfig().(*batchprocessor.Config)
	cfg.SendBatchSize = 8192
	cfg.Timeout = timeout
	cfg.MetadataKeys = make([]string, 0)

	return cfg
}
