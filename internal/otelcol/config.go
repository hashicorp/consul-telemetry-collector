package otelcol

import (
	"context"

	"gopkg.in/yaml.v3"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/otelcol"
	"go.opentelemetry.io/collector/receiver"
)

var cfg string = `
receivers:
  otlp:
    protocols:
      http: {}

exporters:
  logging: {}

service:
  pipelines:
    traces:
      receivers: [otlp]
      exporters: [logging]
`

var othercfg string = `

`

var _ confmap.Provider = (*staticCfg)(nil)

type staticCfg struct {
}

func (c *staticCfg) Retrieve(ctx context.Context, uri string, watcher confmap.WatcherFunc) (*confmap.Retrieved, error) {
	var v string
	switch uri {
	case "static://cfg":
		v = cfg
	case "static://othercfg":
		v = othercfg
	}

	var raw any
	if err := yaml.Unmarshal([]byte(v), &raw); err != nil {
		return nil, err
	}

	return confmap.NewRetrieved(raw)
}

func (c *staticCfg) Scheme() string {
	return "static"
}

func (c *staticCfg) Shutdown(ctx context.Context) error {
	return nil
}

//	Retrieve(ctx context.Context, uri string, watcher WatcherFunc) (*Retrieved, error)
//
//	// Scheme returns the location scheme used by Retrieve.
//	Scheme() string
//
//	// Shutdown signals that the configuration for which this Provider was used to
//	// retrieve values is no longer in use and the Provider should close and release
//	// any resources that it may have created.
//	//
//	// This method must be called when the Collector service ends, either in case of
//	// success or error. Retrieve cannot be called after Shutdown.
//	//
//	// Should never be called concurrently with itself or with Retrieve.
//	// If ctx is cancelled should return immediately with an error.
//	Shutdown(ctx context.Context) error

func ReceiverConfig(f map[component.Type]receiver.Factory) map[component.Type]component.Config {
	cfg := make(map[component.Type]component.Config)
	for t, component := range f {
		cfg[t] = component.CreateDefaultConfig()
	}

	return cfg
}

func newStaticProvider() confmap.Provider {
	return &staticCfg{}
}

func Provider() (otelcol.ConfigProvider, error) {
	return otelcol.NewConfigProvider(otelcol.ConfigProviderSettings{
		ResolverSettings: confmap.ResolverSettings{
			URIs:       []string{"static://cfg", "static://othercfg"},
			Providers:  map[string]confmap.Provider{"static": newStaticProvider()},
			Converters: []confmap.Converter{},
		},
	})
}
