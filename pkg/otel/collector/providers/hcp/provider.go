package hcp

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"

	"github.com/hashicorp/hcp-sdk-go/resource"

	"github.com/hashicorp/consul-telemetry-collector/internal/hcp"
	confresolver "github.com/hashicorp/consul-telemetry-collector/pkg/otel/config"
	"github.com/hashicorp/consul-telemetry-collector/pkg/otel/config/helpers/exporters"
)

type hcpProvider struct {
	otlpHTTPEndpoint string
	client           hcp.TelemetryClient
	clientID         string
	clientSecret     string
	shutdownCh       chan struct{}
}

const scheme = "hcp"
const schemePrefix = scheme + ":"

var _ confmap.Provider = (*hcpProvider)(nil)

// NewProvider creates a new static in memory configmap provider
func NewProvider(
	forwarderEndpoint string,
	client hcp.TelemetryClient,
	clientID,
	clientSecret string,
) confmap.Provider {
	return &hcpProvider{
		otlpHTTPEndpoint: forwarderEndpoint,
		client:           client,
		clientID:         clientID,
		clientSecret:     clientSecret,
		shutdownCh:       make(chan struct{}),
	}
}

func (m *hcpProvider) Retrieve(
	ctx context.Context,
	uri string,
	change confmap.WatcherFunc,
) (*confmap.Retrieved, error) {
	if !strings.HasPrefix(uri, m.Scheme()) {
		return nil, fmt.Errorf("%q uri is not supported by %q provider", uri, m.Scheme())
	}

	_, err := resource.FromString(strings.TrimLeft(uri, schemePrefix))
	if err != nil {
		return nil, fmt.Errorf("unable to parse %q uri as HCP resource URL %w", uri, err)
	}

	c, intermediateCfg, err := confresolver.DefaultConfig(
		&confresolver.DefaultParams{
			OtlpHTTPEndpoint: m.otlpHTTPEndpoint,
			Client:           m.client,
			ClientID:         m.clientID,
			ClientSecret:     m.clientSecret,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failure building configuration %w", err)
	}

	// Start setup for the service
	c.Service.Telemetry = confresolver.Telemetry()
	c.Service.Extensions = intermediateCfg.Extensions

	// Start setup for our different pipelines
	// Inmem is going to filter out the HCP exporter
	inmemPipelineCfg := intermediateCfg.
		Clone().
		FilterExporter(exporters.HCPExporterID).
		ToPipelineConfig()

	inmemID := component.NewID(component.DataTypeMetrics)
	c.Service.Pipelines[inmemID] = inmemPipelineCfg

	hcpPipelineCfg := intermediateCfg.
		Clone().
		FilterExporter(exporters.BaseOtlpExporterID).
		ToPipelineConfig()
	hcpID := component.NewIDWithName(component.DataTypeMetrics, "hcp")
	c.Service.Pipelines[hcpID] = hcpPipelineCfg

	go func() {
		ticker := time.NewTicker(time.Minute)
		for {
			select {
			case <-ctx.Done():
			case <-m.shutdownCh:
				return
			case <-ticker.C:
				if m.configChange() {
					change(&confmap.ChangeEvent{})
					return
				}
			}
		}
	}()

	conf := confmap.New()
	err = conf.Marshal(c)
	if err != nil {
		return nil, err
	}
	return confmap.NewRetrieved(conf.ToStringMap())
}

func (m *hcpProvider) Scheme() string {
	return "hcp"
}

func (m *hcpProvider) Shutdown(ctx context.Context) error {
	close(m.shutdownCh)
	return nil
}

func (m *hcpProvider) configChange() bool {
	// changed := m.client.configChange()
	// return changed
	return false
}
