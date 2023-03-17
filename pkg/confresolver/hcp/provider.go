package hcp

import (
	"context"
	"fmt"
	"strings"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"

	"github.com/hashicorp/consul-telemetry-collector/internal/hcp"
	"github.com/hashicorp/consul-telemetry-collector/pkg/confresolver"
	"github.com/hashicorp/consul-telemetry-collector/pkg/confresolver/confhelper"
)

type hcpProvider struct {
	otlpHTTPEndpoint string
	client           *hcp.Client
}

const scheme = "hcp"
const schemePrefix = scheme + ":"

var _ confmap.Provider = (*hcpProvider)(nil)

// NewProvider creates a new static in memory configmap provider
func NewProvider(forwarderEndpoint string) confmap.Provider {
	return &hcpProvider{
		otlpHTTPEndpoint: forwarderEndpoint,
	}
}

func (m *hcpProvider) Retrieve(ctx context.Context, uri string, change confmap.WatcherFunc) (*confmap.Retrieved,
	error) {
	if !strings.HasPrefix(uri, m.Scheme()) {
		return nil, fmt.Errorf("%q uri is not supported by %q provider", uri, m.Scheme())
	}

	// resource, err := resource.FromString(strings.TrimLeft(uri, schemePrefix))
	// if err != nil {
	// 	return nil, fmt.Errorf("unable to parse %q uri as HCP resource URL %w", uri, err)
	// }

	// _ = resource

	c := &confresolver.Config{}
	pipeline := c.NewPipeline(component.DataTypeMetrics)
	hcpPipeline := c.NewPipelineWithName(component.DataTypeMetrics, "hcp")

	// receivers
	confhelper.OTLPReceiver(c, pipeline, hcpPipeline)

	// processors
	confhelper.MemoryLimiter(c, pipeline, hcpPipeline)

	// put other processors here
	// follow recommended practices: https://github.com/open-telemetry/opentelemetry-collector/tree/main/processor#recommended-processors

	// get filtered from hcp
	// c.NewProcessor(component.NewID("filterprocessor", hcpPipeline)

	c.NewProcessor(component.NewID("batch"), pipeline, hcpPipeline)

	confhelper.Ballast(c)

	c.Service.Telemetry = confresolver.Telemetry()

	// Set oauth2client extension
	// oauth2auth := c.NewExtensions(component.NewIDWithName("oauth2client", "hcp"))
	// oauth2auth.Set("issuer_url", "https://auth.idp.hashicorp.cloud")
	// oauth2auth.Set("audience", "https://api.hashicorp.com")
	// oauth2auth.Set("client_id", "")
	// oauth2auth.Set("client_secret", "")
	// oauth2auth.Set("token_url", "")
	// endpointParams := oauth2auth.SetMap("endpoint_params")
	// endpointParams.Set("audience", "")

	// fetch otlp endpoint from the HCP client here
	// telemetryEndpoint := m.client.Telemetry()

	telemetryEndpoint := "https://hcp-metrics-endpoint"
	c.NewExporter(component.NewID("logging"), pipeline, hcpPipeline)
	otlphttp := c.NewExporter(component.NewIDWithName("otlphttp", "hcp"), hcpPipeline)
	otlphttp.Set("endpoint", telemetryEndpoint)
	// otlphttp.SetMap("auth").Set("authenticator", component.NewID("oidc/hcp"))

	if m.otlpHTTPEndpoint != "" {
		c.PushExporterOnPipeline(pipeline, component.NewID("otlphttp"))
	}

	changeCh := m.configChange()
	go func() {
		for {
			select {
			case <-ctx.Done():
			case <-changeCh:
				change(&confmap.ChangeEvent{})
			}
		}

	}()

	conf := confmap.New()
	err := conf.Marshal(c)
	if err != nil {
		return nil, err
	}
	return confmap.NewRetrieved(conf.ToStringMap())
}

func (m *hcpProvider) Scheme() string {
	return "hcp"
}

func (m *hcpProvider) Shutdown(ctx context.Context) error {
	return nil
}

func (m *hcpProvider) configChange() <-chan struct{} {
	changeCh := make(chan struct{})
	go func() {
		// m.client.configChange
		changeCh <- struct{}{}
		close(changeCh)
	}()
	return changeCh
}