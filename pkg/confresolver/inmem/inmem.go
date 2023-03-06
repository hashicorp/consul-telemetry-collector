package inmem

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"

	"github.com/hashicorp/consul-telemetry-collector/pkg/confresolver"
)

type memProvider struct{}

var _ confmap.Provider = (*memProvider)(nil)

// NewProvider creates a new static in memory configmap provider
func NewProvider() confmap.Provider {
	return new(memProvider)
}

func (m *memProvider) Retrieve(_ context.Context, _ string, _ confmap.WatcherFunc) (*confmap.Retrieved,
	error) {

	c := &confresolver.Config{}
	pipeline := c.NewPipeline(component.DataTypeTraces)
	receiver := c.NewReceiver(pipeline, component.NewID("otlp"))
	receiver.Map("protocols").Map("http")
	c.NewExporter(pipeline, component.NewID("logging"))
	c.Service.Telemetry = confresolver.Telemetry()

	conf := confmap.New()
	err := conf.Marshal(c)
	if err != nil {
		return nil, err
	}
	return confmap.NewRetrieved(conf.ToStringMap())
}

func (m *memProvider) Scheme() string {
	return "mem"
}

func (m *memProvider) Shutdown(ctx context.Context) error {
	return nil
}
