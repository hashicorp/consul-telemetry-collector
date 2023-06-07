// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/service"

	"github.com/hashicorp/consul-telemetry-collector/internal/hcp"
	"github.com/hashicorp/consul-telemetry-collector/internal/otel/config/helpers/exporters"
	"github.com/hashicorp/consul-telemetry-collector/internal/otel/config/helpers/extensions"
	"github.com/hashicorp/consul-telemetry-collector/internal/otel/config/helpers/processors"
	"github.com/hashicorp/consul-telemetry-collector/internal/otel/config/helpers/receivers"
)

// Params are the inputs to the configuration building process. Only some config requires
// these inputs.
type Params struct {
	OtlpHTTPEndpoint string
	Client           hcp.TelemetryClient
	ClientID         string
	ClientSecret     string
	ResourceID       string
}

// PipelineConfigBuilder defines a basic list of pipeline component IDs for a service.PipelineConfig.
func PipelineConfigBuilder(p *Params) service.PipelineConfig {
	baseCfg := service.PipelineConfig{
		Processors: []component.ID{
			processors.MemoryLimiterID,
			// add your processors here
			processors.BatchProcessorID,
		},
		Receivers: []component.ID{receivers.EnvoyReceiverID, receivers.PrometheusReceiverID},
		Exporters: []component.ID{
			exporters.LoggingExporterID,
		},
	}

	includeHCPPipeline := p.ClientID != "" && p.ClientSecret != "" && p.Client != nil
	if includeHCPPipeline {
		baseCfg.Exporters = append(baseCfg.Exporters, exporters.HCPExporterID)
	} else if p.OtlpHTTPEndpoint != "" {
		baseCfg.Exporters = append(baseCfg.Exporters, exporters.BaseOtlpExporterID)
	}

	return baseCfg
}

// Opts is a variadic type passed in as a way  of manipulating a list of components.
type Opts func([]component.ID) []component.ID

// WithExtOauthClientID is an Opt function to add the oauthclient id to the list of extensions
// NOTE: this extension will require.
func WithExtOauthClientID(ext []component.ID) []component.ID {
	return append(ext, extensions.OauthClientID)
}

// WithFilterProcessor is an Opt function to add the filter processor to a list of processors.
func WithFilterProcessor(procesors []component.ID) []component.ID {
	return append(procesors, processors.FilterProcessorID)
}

// WithResourceProcessor adds the resource processor to a list of processors. It should go after the filter processor to ensure that we do not operate on signals that we won't forward.
func WithResourceProcessor(prcs []component.ID) []component.ID {
	return append(prcs, processors.ResourceProcessorID)
}

// ExtensionBuilder builds a list of extension IDs. Optionally we can include more ids with variadic opts.
func ExtensionBuilder(opts ...Opts) []component.ID {
	base := []component.ID{
		extensions.BallastID,
	}
	for _, opt := range opts {
		base = opt(base)
	}
	return base
}

// ProcessorBuilder returns a list of processor IDs.
// The provided IDs inserted between the memory limiter and batch processor.
func ProcessorBuilder(opts ...Opts) []component.ID {
	base := []component.ID{
		processors.MemoryLimiterID,
	}
	for _, opt := range opts {
		base = opt(base)
	}

	base = append(base, processors.BatchProcessorID)
	return base
}
