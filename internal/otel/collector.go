// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package otel

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/featuregate"
	"go.opentelemetry.io/collector/otelcol"

	"github.com/hashicorp/consul-telemetry-collector/internal/hcp"
	"github.com/hashicorp/consul-telemetry-collector/internal/otel/config/helpers/exporters"
	"github.com/hashicorp/consul-telemetry-collector/internal/version"
)

// Collector is an interface that is satisfied by the otelcol.Collector struct.
// This allows us to wrap the opentelemetry collector and not necessarily run it ourselves.
type Collector interface {
	Run(context.Context) error
	GetState() otelcol.State
	Shutdown()
}

// CollectorCfg is the configuration needed to start the collector.
type CollectorCfg struct {
	ClientID          string
	ClientSecret      string
	ResourceID        string
	Client            hcp.TelemetryClient
	ForwarderEndpoint string
	ExporterConfig    *OTLPExporterConfig
}

// OTLPExporterConfig holds the type and
type OTLPExporterConfig struct {
	ExporterConfig exporters.ExporterConfig
	Type           string
}

const otelFeatureGate = "telemetry.useOtelForInternalMetrics"

// NewCollector will create a new open-telemetry collector service and configuration based on the provided values.
func NewCollector(cfg CollectorCfg) (Collector, error) {
	// enable otel for collector internal metrics
	if err := featuregate.GlobalRegistry().Set(otelFeatureGate, true); err != nil {
		return nil, err
	}

	factories, err := components()
	if err != nil {
		return nil, err
	}

	provider, err := newProvider(cfg)
	if err != nil {
		return nil, err
	}

	set := otelcol.CollectorSettings{
		Factories: factories,
		BuildInfo: component.BuildInfo{
			Command:     "consul-telemetry-collector",
			Description: "consul-telemetry-collector is a Consul specific build of the open-telemetry collector",
			Version:     version.Version,
		},
		DisableGracefulShutdown: true,
		ConfigProvider:          provider,
		LoggingOptions:          nil,
		SkipSettingGRPCLogger:   false,
	}

	return otelcol.NewCollector(set)
}
