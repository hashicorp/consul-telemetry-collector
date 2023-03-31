package otel

import (
	"fmt"

	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/otelcol"

	internalhcp "github.com/hashicorp/consul-telemetry-collector/internal/hcp"
	"github.com/hashicorp/consul-telemetry-collector/pkg/otel/providers/external"
	"github.com/hashicorp/consul-telemetry-collector/pkg/otel/providers/hcp"
)

type collectorCfg struct {
	clientID     string
	clientSecret string
	resourceID   string
	client       internalhcp.TelemetryClient
}

type collectorOpts func(*collectorCfg)

// revive:disable:unexported-return

// WithCloud adds cloud parameters to the collector creation
func WithCloud(resourceID string, clientID string, clientSecret string, client internalhcp.TelemetryClient) collectorOpts {
	return func(o *collectorCfg) {
		o.resourceID = resourceID
		o.clientID = clientID
		o.clientSecret = clientSecret
		o.client = client
	}
}

//revive:enable:unexported-return

func newProvider(
	forwarderEndpoint string,
	resourceID string,
	clientID string,
	clientSecret string,
	client internalhcp.TelemetryClient,
) (otelcol.ConfigProvider, error) {
	uris := []string{"external:"}
	if resourceID != "" {
		uris = append(uris, fmt.Sprintf("hcp:%s", resourceID))
	}
	resolver := confmap.ResolverSettings{
		URIs: uris,
		Providers: makeMapProvidersMap(
			external.NewProvider(forwarderEndpoint),
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
