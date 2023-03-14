package otelcol

import (
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/otelcol"

	"github.com/hashicorp/consul-telemetry-collector/pkg/confresolver/inmem"
)

func newConfigProvider(forwarderEndpoint string) (otelcol.ConfigProvider, error) {
	resolver := confmap.ResolverSettings{
		URIs: []string{"inmem://"},
		Providers: map[string]confmap.Provider{
			"inmem": inmem.NewProvider(forwarderEndpoint),
		},
		Converters: []confmap.Converter{},
	}

	return otelcol.NewConfigProvider(otelcol.ConfigProviderSettings{
		ResolverSettings: resolver,
	})
}
