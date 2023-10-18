// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package external

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"

	"github.com/hashicorp/consul-telemetry-collector/internal/otel/config"
)

type externalProvider struct {
	exporterID component.ID
	exporter   config.Exporter
}

var _ confmap.Provider = (*externalProvider)(nil)

// NewProvider creates a new static in memory configmap provider.
func NewProvider(exporterID component.ID, exporter config.Exporter) confmap.Provider {
	return &externalProvider{
		exporterID: exporterID,
		exporter:   exporter,
	}
}

func (m *externalProvider) Retrieve(_ context.Context, _ string, _ confmap.WatcherFunc) (*confmap.Retrieved,
	error) {
	// Create new empty configuration
	c := config.NewConfig()

	// 1. Setup Extensions
	c.Service.Telemetry = config.Telemetry()

	// 2. Setup Extensions
	extensions := config.ExtensionBuilder()
	// so far we don't use the params in building extensions for the external provider
	err := c.EnrichWithExtensions(extensions, nil)
	if err != nil {
		return nil, err
	}

	// 3. Build external pipeline
	externalParams := &config.Params{
		OtlpHTTPEndpoint: m.otlpHTTPEndpoint,
	}
	externalCfg := config.PipelineConfigBuilder(externalParams)
	externalID := component.NewID(component.DataTypeMetrics)
	err = c.EnrichWithPipelineCfg(externalCfg, externalParams, externalID)
	if err != nil {
		return nil, fmt.Errorf("failed to add config to pipeline. provider:external, err: %w", err)
	}

	conf := confmap.New()
	err = conf.Marshal(c)
	if err != nil {
		return nil, err
	}
	return confmap.NewRetrieved(conf.ToStringMap())
}

func (m *externalProvider) Scheme() string {
	return "external"
}

func (m *externalProvider) Shutdown(_ context.Context) error {
	return nil
}
