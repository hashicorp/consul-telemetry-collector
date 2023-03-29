package collector

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/featuregate"
	"go.opentelemetry.io/collector/otelcol"

	"github.com/hashicorp/consul-telemetry-collector/pkg/version"
)

const otelFeatureGate = "telemetry.useOtelForInternalMetrics"

// New will create a new open-telemetry collector service and configuration based on the provided values
func New(ctx context.Context, forwarderEndpoint string, opts ...collectorOpts) (Collector,
	error) {
	cfg := &collectorCfg{}
	for _, opt := range opts {
		opt(cfg)
	}
	// enable otel for collector internal metrics
	if err := featuregate.GlobalRegistry().Set(otelFeatureGate, true); err != nil {
		return nil, err
	}

	factories, err := components()
	if err != nil {
		return nil, err
	}

	provider, err := newProvider(forwarderEndpoint, cfg.resourceID, cfg.clientID, cfg.clientSecret, cfg.client)
	if err != nil {
		return nil, err
	}

	set := otelcol.CollectorSettings{
		Factories: factories,
		BuildInfo: component.BuildInfo{
			Command:     "consul-telemetry-collector",
			Description: "consul-telemetry-collector is a Consul specific build of the open-telemetry collector",
			Version:     version.Version,
		},
		DisableGracefulShutdown: true,
		ConfigProvider:          provider,
		LoggingOptions:          nil,
		SkipSettingGRPCLogger:   false,
	}

	return otelcol.NewCollector(set)
}