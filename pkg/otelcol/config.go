package otelcol

import (
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/otelcol"

	"github.com/hashicorp/consul-telemetry-collector/pkg/confresolver/inmem"
)

func newConfigProvider() (otelcol.ConfigProvider, error) {
	return otelcol.NewConfigProvider(otelcol.ConfigProviderSettings{
		ResolverSettings: confmap.ResolverSettings{
			URIs: []string{"mem://"},
			Providers: map[string]confmap.Provider{
				"mem": inmem.NewProvider(),
			},
			Converters: []confmap.Converter{},
		},
	})
}
