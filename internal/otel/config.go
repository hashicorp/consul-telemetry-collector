// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package otel

import (
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/otelcol"

	"github.com/hashicorp/consul-telemetry-collector/internal/otel/config"
	"github.com/hashicorp/consul-telemetry-collector/internal/otel/providers/external"
	"github.com/hashicorp/consul-telemetry-collector/internal/otel/providers/hcp"
)

func newProvider(cfg CollectorCfg) (otelcol.ConfigProvider, error) {
	uris := []string{"external:"}
	if cfg.ResourceID != "" {
		uris = append(uris, fmt.Sprintf("hcp:%s", cfg.ResourceID))
	}

	var exportID component.ID
	var exportConfig config.Exporter
	// TODO: Make sure a nil ExporterConfig is tested
	if cfg.ExporterConfig != nil {
		exportID = cfg.ExporterConfig.ID
		exportConfig = cfg.ExporterConfig.Exporter
	}

	resolver := confmap.ResolverSettings{
		URIs: uris,
		Providers: makeMapProvidersMap(
			external.NewProvider(exportID, exportConfig),
			hcp.NewProvider(cfg.ExporterConfig, cfg.Client, cfg.ClientID, cfg.ClientSecret),
		),
		Converters: []confmap.Converter{},
	}

	return otelcol.NewConfigProvider(otelcol.ConfigProviderSettings{
		ResolverSettings: resolver,
	})
}

func makeMapProvidersMap(providers ...confmap.Provider) map[string]confmap.Provider {
	ret := make(map[string]confmap.Provider, len(providers))
	for _, provider := range providers {
		ret[provider.Scheme()] = provider
	}
	return ret
}
