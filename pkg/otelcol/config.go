package otelcol

import (
	"fmt"

	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/otelcol"

	internalhcp "github.com/hashicorp/consul-telemetry-collector/internal/hcp"
	"github.com/hashicorp/consul-telemetry-collector/pkg/confresolver/hcp"
	"github.com/hashicorp/consul-telemetry-collector/pkg/confresolver/inmem"
)

func newConfigProvider(forwarderEndpoint string, resourceURL string,
	clientID string, clientSecret string, client internalhcp.TelemetryClient) (otelcol.ConfigProvider, error) {
	uris := []string{"inmem:"}
	if resourceURL != "" {
		uris = append(uris, fmt.Sprintf("hcp:%s", resourceURL))
	}
	resolver := confmap.ResolverSettings{
		URIs: uris,
		Providers: makeMapProvidersMap(
			inmem.NewProvider(forwarderEndpoint),
			hcp.NewProvider(forwarderEndpoint, client, clientID, clientSecret),
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
