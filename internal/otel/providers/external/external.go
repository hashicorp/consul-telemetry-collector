// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package external

import (
	"context"
	"fmt"
	"os"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"gopkg.in/yaml.v2"

	"github.com/hashicorp/consul-telemetry-collector/internal/otel/config"
)

type externalProvider struct {
	otlpHTTPEndpoint string
	overridesPath    string
}

var _ confmap.Provider = (*externalProvider)(nil)

// NewProvider creates a new static in memory configmap provider.
func NewProvider(forwarderEndpoint, overridesFP string) confmap.Provider {
	return &externalProvider{
		otlpHTTPEndpoint: forwarderEndpoint,
		overridesPath:    overridesFP,
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

	if m.overridesPath != "" {
		overrideCfg, err := loadOveride(m.overridesPath)
		if err != nil {
			return nil, fmt.Errorf("failure to parse overrides")
		}

		if err := conf.Merge(overrideCfg); err != nil {
			return nil, err
		}
	}

	return confmap.NewRetrieved(conf.ToStringMap())
}

func loadOveride(filePath string) (*confmap.Conf, error) {
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse external config: %w", err)
	}
	raw := map[string]interface{}{}
	err = yaml.Unmarshal(yamlFile, raw)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal external config: %w", err)
	}
	overrideCfg := confmap.New()
	err = overrideCfg.Marshal(raw)
	if err != nil {
		return nil, err
	}
	return overrideCfg, nil
}

func (m *externalProvider) Scheme() string {
	return "external"
}

func (m *externalProvider) Shutdown(_ context.Context) error {
	return nil
}
