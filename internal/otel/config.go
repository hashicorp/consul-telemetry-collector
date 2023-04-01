package otel

import (
	"fmt"

	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/otelcol"

	"github.com/hashicorp/consul-telemetry-collector/internal/otel/providers/external"
	"github.com/hashicorp/consul-telemetry-collector/internal/otel/providers/hcp"
)

// revive:disable:unexported-return

//revive:enable:unexported-return

func newProvider(
	cfg CollectorCfg,
) (otelcol.ConfigProvider, error) {
	uris := []string{"external:"}
	if cfg.ResourceID != "" {
		uris = append(uris, fmt.Sprintf("hcp:%s", cfg.ResourceID))
	}
	resolver := confmap.ResolverSettings{
		URIs: uris,
		Providers: makeMapProvidersMap(
			external.NewProvider(cfg.ForwarderEndpoint),
			hcp.NewProvider(cfg.ForwarderEndpoint, cfg.Client, cfg.ClientID, cfg.ClientSecret),
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
