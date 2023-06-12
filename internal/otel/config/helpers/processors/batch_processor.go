// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package processors

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/processor/batchprocessor"
)

const batchProcessorName = "batch"

// BatchProcessorID is the component id of the batch processor.
var BatchProcessorID component.ID = component.NewID(batchProcessorName)

// BatchProcessorCfg  generates the config for a batch processor.
func BatchProcessorCfg() *batchprocessor.Config {
	factory := batchprocessor.NewFactory()
	cfg := factory.CreateDefaultConfig().(*batchprocessor.Config)

	return cfg
}
