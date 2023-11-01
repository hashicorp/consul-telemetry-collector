// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package otel

import (
	"fmt"

	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/otelcol"

	"github.com/hashicorp/consul-telemetry-collector/internal/otel/providers"
	"github.com/hashicorp/consul-telemetry-collector/internal/otel/providers/external"
	"github.com/hashicorp/consul-telemetry-collector/internal/otel/providers/hcp"
)

func newProvider(cfg CollectorCfg) (otelcol.ConfigProvider, error) {
	uris := []string{"external:"}
	if cfg.ResourceID != "" {
		uris = append(uris, fmt.Sprintf("hcp:%s", cfg.ResourceID))
	}

	params := providers.SharedParams{
		BatchTimeout: cfg.BatchTimeout,
		MetricsPort:  cfg.MetricsPort,
		EnvoyPort:    cfg.EnvoyPort,
	}

	resolver := confmap.ResolverSettings{
		URIs: uris,
		Providers: makeMapProvidersMap(
			external.NewProvider(cfg.ExporterConfig, params),
			hcp.NewProvider(cfg.ExporterConfig, cfg.Client, cfg.ClientID, cfg.ClientSecret, params),
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
