package inmem

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"

	"github.com/hashicorp/consul-telemetry-collector/pkg/confresolver"
	"github.com/hashicorp/consul-telemetry-collector/pkg/confresolver/confhelper"
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

	c := &confresolver.Config{}
	pipeline := c.NewPipeline(component.DataTypeMetrics)
	confhelper.OTLPReceiver(c, pipeline)

	c.NewExporter(component.NewID("logging"), pipeline)

	confhelper.MemoryLimiter(c, pipeline)

	// put other processors here
	// follow recommended practices: https://github.com/open-telemetry/opentelemetry-collector/tree/main/processor#recommended-processors

	c.NewProcessor(component.NewID("batch"), pipeline)

	confhelper.Ballast(c)

	c.Service.Telemetry = confresolver.Telemetry()

	if m.otlpHTTPEndpoint != "" {
		otlphttp := c.NewExporter(component.NewID("otlphttp"), pipeline)
		otlphttp.Set("endpoint", m.otlpHTTPEndpoint)
	}

	conf := confmap.New()
	err := conf.Marshal(c)
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
