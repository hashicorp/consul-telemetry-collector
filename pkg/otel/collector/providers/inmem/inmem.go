package inmem

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"

	"github.com/hashicorp/consul-telemetry-collector/pkg/otel/config"
)

type inmemProvider struct {
	otlpHTTPEndpoint string
}

var _ confmap.Provider = (*inmemProvider)(nil)

// NewProvider creates a new static in memory configmap provider
func NewProvider(forwarderEndpoint string) confmap.Provider {
	return &inmemProvider{
		otlpHTTPEndpoint: forwarderEndpoint,
	}
}

func (m *inmemProvider) Retrieve(_ context.Context, _ string, _ confmap.WatcherFunc) (*confmap.Retrieved,
	error) {
	// Overall configuration that will hold all receivers/exporters/processors/connectors/extensions
	// and service config
	c, intermediateCfg, err := config.DefaultConfig(
		&config.DefaultParams{
			OtlpHTTPEndpoint: m.otlpHTTPEndpoint,
		},
	)
	if err != nil {
		return nil, err
	}

	// Start Setting up the service
	c.Service.Telemetry = config.Telemetry()
	c.Service.Extensions = intermediateCfg.Extensions

	// Start setup for our different pipelines
	pipelineConfig := intermediateCfg.ToPipelineConfig()
	inmemID := component.NewID(component.DataTypeMetrics)
	c.Service.Pipelines[inmemID] = pipelineConfig

	// pipelineConfig2 := intermediateCfg.Clone().FilterExporter(&memoryLimiterId).ToPipelineConfig()
	// hcpID := component.NewIDWithName(component.DataTypeMetrics, "hcp")
	// c.Service.Pipelines[hcpID] = pipelineConfig2

	conf := confmap.New()
	err = conf.Marshal(c)
	if err != nil {
		return nil, err
	}
	return confmap.NewRetrieved(conf.ToStringMap())
}

func (m *inmemProvider) Scheme() string {
	return "inmem"
}

func (m *inmemProvider) Shutdown(ctx context.Context) error {
	return nil
}
