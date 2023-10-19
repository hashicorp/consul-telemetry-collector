// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package hcp

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"

	"github.com/hashicorp/consul-telemetry-collector/internal/hcp"
	"github.com/hashicorp/consul-telemetry-collector/internal/otel/config"
	"github.com/hashicorp/hcp-sdk-go/resource"
)

type hcpProvider struct {
	exporterConfig *config.ExportConfig
	client         hcp.TelemetryClient
	clientID       string
	clientSecret   string
	shutdownCh     chan struct{}
}

const scheme = "hcp"
const schemePrefix = scheme + ":"

var _ confmap.Provider = (*hcpProvider)(nil)

// NewProvider creates a new static in memory configmap provider.
func NewProvider(
	exporterConfig *config.ExportConfig,
	client hcp.TelemetryClient,
	clientID,
	clientSecret string,
) confmap.Provider {
	return &hcpProvider{
		exporterConfig: exporterConfig,
		client:         client,
		clientID:       clientID,
		clientSecret:   clientSecret,
		shutdownCh:     make(chan struct{}),
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

	r, err := resource.FromString(strings.TrimLeft(uri, schemePrefix))
	if err != nil {
		return nil, fmt.Errorf("unable to parse %q uri as HCP resource URL %w", uri, err)
	}

	// Create new empty configuration
	c := config.NewConfig()

	// 1. Setup Telemetery
	c.Service.Telemetry = config.Telemetry()

	// 2. Setup Extensions
	extensions := config.ExtensionBuilder(config.WithExtOauthClientID)
	// in this set of extension IDs we want the WithExtOauthClientID which requires the params to build
	// the actual extension.
	hcpParams := &config.Params{
		ExporterConfig: m.exporterConfig,
		Client:         m.client,
		ClientID:       m.clientID,
		ClientSecret:   m.clientSecret,
		ResourceID:     r.String(),
	}
	err = c.EnrichWithExtensions(extensions, hcpParams)
	if err != nil {
		return nil, err
	}

	// 3. Build pipeline configurations and enrich the config with them
	// 3. A: Build HCP pipeline
	hcpPipelineCfg := config.PipelineConfigBuilder(hcpParams)

	// Set the filter processor on the config
	hcpPipelineCfg.Processors = config.ProcessorBuilder(config.WithFilterProcessor, config.WithResourceProcessor)

	hcpID := component.NewIDWithName(component.DataTypeMetrics, "hcp")
	err = c.EnrichWithPipelineCfg(hcpPipelineCfg, hcpParams, hcpID)
	if err != nil {
		return nil, err
	}

	// 3. B: Build external pipeline We need to build this external pipeline because of how otel merges configuration.
	// It does _not_ perform a deep merge and instead performs an overriding merge at the highest matching level. This
	// behavior means that the service object will never match between the HCP provider and the External provider
	// without ensuring that the external generator _also_ executes for HCP. The external generator will ensure that the
	// forward component.ID and other components are included in the service stanza and are activated by the collector.
	// An improvement here would be to separate the service stanza creation from the HCP or External generators. This
	// would allow component configuration to happen separately from the service stanza and removing repeated work.
	externalParams := &config.Params{
		ExporterConfig: m.exporterConfig,
	}
	externalCfg := config.PipelineConfigBuilder(externalParams)
	externalID := component.NewID(component.DataTypeMetrics)
	err = c.EnrichWithPipelineCfg(externalCfg, externalParams, externalID)
	if err != nil {
		return nil, err
	}

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

func (m *hcpProvider) Shutdown(_ context.Context) error {
	close(m.shutdownCh)
	return nil
}

func (m *hcpProvider) configChange() bool {
	return false
}
