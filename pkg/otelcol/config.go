package otelcol

import (
	"fmt"

	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/otelcol"

	"github.com/hashicorp/consul-telemetry-collector/pkg/confresolver/hcp"
	"github.com/hashicorp/consul-telemetry-collector/pkg/confresolver/inmem"
)

func newConfigProvider(forwarderEndpoint string, resourceURL string) (otelcol.ConfigProvider, error) {
	uris := []string{"inmem:"}
	if resourceURL != "" {
		uris = append(uris, fmt.Sprintf("hcp:%s", resourceURL))
	}
	resolver := confmap.ResolverSettings{
		URIs: uris,
		Providers: map[string]confmap.Provider{
			"inmem": inmem.NewProvider(forwarderEndpoint),
			"hcp":   hcp.NewProvider(forwarderEndpoint),
		},
		Converters: []confmap.Converter{},
	}

	return otelcol.NewConfigProvider(otelcol.ConfigProviderSettings{
		ResolverSettings: resolver,
	})
}
