package config

import (
	"github.com/hashicorp/consul-telemetry-collector/internal/hcp"
	"github.com/hashicorp/consul-telemetry-collector/pkg/otel/config/helpers/exporters"
	"github.com/hashicorp/consul-telemetry-collector/pkg/otel/config/helpers/extensions"
	"github.com/hashicorp/consul-telemetry-collector/pkg/otel/config/helpers/processors"
	"github.com/hashicorp/consul-telemetry-collector/pkg/otel/config/helpers/receivers"
	oauth "github.com/open-telemetry/opentelemetry-collector-contrib/extension/oauth2clientauthextension"
	"go.opentelemetry.io/collector/component"
)

// DefaultParams is the default parameters passed to the default config builder below.
type DefaultParams struct {
	OtlpHTTPEndpoint string
	Client           hcp.TelemetryClient
	ClientID         string
	ClientSecret     string
}

// DefaultConfig generates a default config for pipelines. This will likely change
// so that pipelines define their components and then we generate the config (the one built in intermediate)
// from their ids vs doing the base configuration and then filtering out the ids for the specific pipeline.
func DefaultConfig(p *DefaultParams) (*Config, *IntermediateConfig, error) {
	includeHCPPipeline := p.ClientID != "" && p.ClientSecret != "" && p.Client != nil
	// Overall configuration that will hold all receivers/exporters/processors/connectors/extensions
	// and service config
	c := NewConfig()

	intermediateCfg := NewIntermediateConfig()

	// Setup the otlp receiver
	otlpID, otlpConfig := receivers.OtlpReceiverCfg()
	c.Receivers[otlpID] = otlpConfig
	intermediateCfg.Receivers = append(intermediateCfg.Receivers, otlpID)

	// Setup the memory limiter
	memoryLimiterID, memoryLimiterCfg := processors.MemoryLimiterCfg()
	c.Processors[memoryLimiterID] = memoryLimiterCfg
	intermediateCfg.Processors = append(intermediateCfg.Processors, memoryLimiterID)

	//  Add your processors here

	// Setup the batch processor
	batchProcesserID, batchProcessorCfg := processors.BatchProcessorCfg()
	c.Processors[batchProcesserID] = batchProcessorCfg
	intermediateCfg.Processors = append(intermediateCfg.Processors, batchProcesserID)

	// Setup the ballast extension
	ballastExtensionID, ballastExtCfg := extensions.BallastCfg()
	c.Extensions[ballastExtensionID] = ballastExtCfg
	intermediateCfg.Extensions = append(intermediateCfg.Extensions, ballastExtensionID)

	// Setup the logging exporter
	loggingID, loggingConfig := exporters.LogExporterCfg()
	c.Exporters[loggingID] = loggingConfig
	intermediateCfg.Exporters = append(intermediateCfg.Exporters, loggingID)

	// Special oauth extension and exporter for HCP
	var oauthExtensionID component.ID
	var oauthExtensionCfg *oauth.Config
	if includeHCPPipeline {
		// setup oauth extension
		oauthExtensionID, oauthExtensionCfg = extensions.OauthClientCfg(p.ClientID, p.ClientSecret)
		c.Extensions[oauthExtensionID] = oauthExtensionCfg
		intermediateCfg.Extensions = append(intermediateCfg.Extensions, oauthExtensionID)

		metricsEndpoint, err := p.Client.MetricsEndpoint()
		if err != nil {
			return nil, nil, err
		}
		// setup the HCP exporter
		hcpExporterID, hcpExporterCfg := exporters.OtlpExporterHCPCfg(metricsEndpoint, oauthExtensionID)
		c.Exporters[hcpExporterID] = hcpExporterCfg
		intermediateCfg.Exporters = append(intermediateCfg.Exporters, hcpExporterID)
	}

	if p.OtlpHTTPEndpoint != "" {
		// Setup the otlpHTTPEndpoint exporter
		otlpExporterID, otlpExporterCfg := exporters.OtlpExporterCfg(p.OtlpHTTPEndpoint)
		c.Exporters[otlpExporterID] = otlpExporterCfg
		intermediateCfg.Exporters = append(intermediateCfg.Exporters, otlpExporterID)
	}

	return c, intermediateCfg, nil

}
