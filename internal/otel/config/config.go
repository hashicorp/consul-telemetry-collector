package config

import (
	"errors"

	"github.com/hashicorp/consul-telemetry-collector/internal/otel/config/helpers/exporters"
	"github.com/hashicorp/consul-telemetry-collector/internal/otel/config/helpers/extensions"
	"github.com/hashicorp/consul-telemetry-collector/internal/otel/config/helpers/processors"
	"github.com/hashicorp/consul-telemetry-collector/internal/otel/config/helpers/receivers"
	"github.com/hashicorp/go-multierror"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/service"
)

// componentMap is a way of identifying a component and it's specific configuration
type componentMap map[component.ID]any

// Config is a helper type to create a new opentelemetry server configuration.
// It implements a map[string]interface{} representation of the opentelemetry-collector configuration.
// More information can be found here: https://opentelemetry.io/docs/collector/configuration/
type Config struct {
	Receivers  componentMap   `mapstructure:"receivers"`
	Exporters  componentMap   `mapstructure:"exporters"`
	Processors componentMap   `mapstructure:"processors"`
	Connectors componentMap   `mapstructure:"connectors"`
	Extensions componentMap   `mapstructure:"extensions"`
	Service    service.Config `mapstructure:"service"`
}

// NewConfig creates a new config object with all types initialized
func NewConfig() *Config {
	svcConfig := service.Config{}
	svcConfig.Pipelines = make(map[component.ID]*service.PipelineConfig)
	return &Config{
		Receivers:  make(componentMap),
		Exporters:  make(componentMap),
		Processors: make(componentMap),
		Connectors: make(componentMap),
		Extensions: make(componentMap),
		Service:    svcConfig,
	}
}

// EnrichWithPipelineCfg enrichs a Config by taking the IDs specified in a pipeline config
// and builds the corresponding configuration for each component ID. Some of these components
// require a set of params
func (c *Config) EnrichWithPipelineCfg(
	pCfg service.PipelineConfig,
	p *Params,
	pipelineID component.ID,
) error {
	var merr *multierror.Error
	// Receivers
	err := buildComponents(c.Receivers, pCfg.Receivers, p)
	merr = multierror.Append(merr, err)
	// Exporters
	err = buildComponents(c.Exporters, pCfg.Exporters, p)
	merr = multierror.Append(merr, err)
	// Processors
	err = buildComponents(c.Processors, pCfg.Processors, p)
	merr = multierror.Append(merr, err)

	if merr.ErrorOrNil() != nil {
		return merr
	}
	c.Service.Pipelines[pipelineID] = &pCfg
	return nil
}

// EnrichWithExtensions adds the specific configurations for a given list of extension IDs.
// The parameters are sometimes required to build an extension so they should be passed through.
func (c *Config) EnrichWithExtensions(
	extensions []component.ID,
	p *Params,
) error {
	if err := buildComponents(c.Extensions, extensions, p); err != nil {
		return err
	}
	c.Service.Extensions = append(c.Service.Extensions, extensions...)
	return nil
}

// buildComponents takes a componentMap (map[component.ID]any) and a list
// of componentIDs. If the componentMap doesn't yet have that component we
// will build it and attach it to the componentMap for that ID. Otherwise we move
// on
func buildComponents(
	tc componentMap,
	componentIDs []component.ID,
	p *Params,
) error {
	for _, id := range componentIDs {
		if _, ok := tc[id]; !ok {
			component, err := buildComponent(id, p)
			if err != nil {
				return err
			}
			tc[id] = component
		}
	}
	return nil
}

// buildComponent returns a configuration type for a specific ID
func buildComponent(id component.ID, p *Params) (any, error) {
	switch id {
	// receivers
	case receivers.OtlpReceiverID:
		return receivers.OtlpReceiverCfg(), nil
	// processors
	case processors.MemoryLimiterID:
		return processors.MemoryLimiterCfg(), nil
	case processors.BatchProcessorID:
		return processors.BatchProcessorCfg(), nil
	// exporters
	case exporters.LoggingExporterID:
		return exporters.LogExporterCfg(), nil
	case exporters.BaseOtlpExporterID:
		return exporters.OtlpExporterCfg(p.OtlpHTTPEndpoint), nil
	case exporters.HCPExporterID:
		if p.Client == nil {
			return nil, errors.New("parameters must specify a client to build HPC exporter config")
		}
		metricsEndpoint, err := p.Client.MetricsEndpoint()
		if err != nil {
			return nil, err
		}

		return exporters.OtlpExporterHCPCfg(metricsEndpoint, extensions.OauthClientID), nil
	// extensions
	case extensions.BallastID:
		return extensions.BallastCfg(), nil
	case extensions.OauthClientID:
		if p.ClientID == "" || p.ClientSecret == "" {
			return nil, errors.New("parameters must specify a client id and secret to build an Oauth extension")
		}
		return extensions.OauthClientCfg(p.ClientID, p.ClientSecret), nil
	default:
		return nil, errors.New("unsupported component id")
	}
}
