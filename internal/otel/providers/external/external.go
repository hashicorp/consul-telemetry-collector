// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package external

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"

	"github.com/hashicorp/consul-telemetry-collector/internal/otel/config"
)

type externalProvider struct {
	exporterConfig *config.ExporterConfig
	batchTimeout   time.Duration
	metricsPort    int
	envoyPort      int
}

type providerOpts func(*externalProvider)

// WithTestOpts allows us to thread testing functionality into the config builder
func WithTestOpts(metricsPort, envoyPort int, batchTimeout time.Duration) providerOpts {
	return func(hp *externalProvider) {
		hp.batchTimeout = batchTimeout
		hp.metricsPort = metricsPort
		hp.envoyPort = envoyPort
	}
}

var _ confmap.Provider = (*externalProvider)(nil)

// NewProvider creates a new static in memory configmap provider.
func NewProvider(exporterConfig *config.ExporterConfig, opt ...providerOpts) confmap.Provider {
	e := &externalProvider{
		exporterConfig: exporterConfig,
	}

	for _, o := range opt {
		o(e)
	}

	return e
}

func (m *externalProvider) Retrieve(_ context.Context, _ string, _ confmap.WatcherFunc) (*confmap.Retrieved,
	error) {
	// Create new empty configuration
	c := config.NewConfig()

	// 1. Setup Extensions
	c.Service.Telemetry = config.Telemetry(m.metricsPort)

	// 2. Setup Extensions
	extensions := config.ExtensionBuilder()
	// so far we don't use the params in building extensions for the external provider
	err := c.EnrichWithExtensions(extensions, nil)
	if err != nil {
		return nil, err
	}

	// 3. Build external pipeline
	externalParams := &config.Params{
		BatchTimeout:      m.batchTimeout,
		MetricsPort:       m.metricsPort,
		EnvoyListenerPort: m.envoyPort,
	}

	// see if this is an empty component.ID
	if m.exporterConfig != nil {
		externalParams.ExporterConfig = m.exporterConfig
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
