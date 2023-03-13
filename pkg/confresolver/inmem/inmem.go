package inmem

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"

	"github.com/hashicorp/consul-telemetry-collector/pkg/confresolver"
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
	receiver := c.NewReceiver(pipeline, component.NewID("otlp"))
	receiver.SetMap("protocols").SetMap("http")
	c.NewExporter(pipeline, component.NewID("logging"))

	limiter := c.NewProcessor(pipeline, component.NewID("memory_limiter"))
	limiter.Set("check_interval", "1s")
	limiter.Set("limit_percentage", "50")
	limiter.Set("spike_limit_percentage", "30")
	c.NewProcessor(pipeline, component.NewID("batch"))

	c.Service.Telemetry = confresolver.Telemetry()

	if m.otlpHTTPEndpoint != "" {
		otlphttp := c.NewExporter(pipeline, component.NewID("otlphttp"))
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
