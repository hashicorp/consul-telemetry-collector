package exporters

import (
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configauth"

	"github.com/hashicorp/consul-telemetry-collector/internal/version"
)

// otlpHTTPExporterName is the component.ID value used by the otlphttp exporter.
const otlpHTTPExporterName = "otlphttp"
const channelName = "x-channel"
const channelValue = "consul-telemetry-collector"
const resourceIDHeader = "x-hcp-resource-id"

var (
	// HCPExporterID is the id of the HCP otel exporter.
	HCPExporterID = component.NewIDWithName(otlpHTTPExporterName, "hcp")
	// BaseOtlpExporterID is the id of a base otel exporter.
	BaseOtlpExporterID = component.NewID(otlpHTTPExporterName)
)

// ExporterConfig is a base wrapper around the otlphttpexorter which
// we cannot use directly since our golden tests can't handle the comparisons unfortunately.
// https://pkg.go.dev/go.opentelemetry.io/collector/exporter/otlphttpexporter@v0.72.0#section-readme
type ExporterConfig struct {
	// Endpoint to send strings to it
	Endpoint string `mapstructure:"endpoint"`
	// Auth configuration for the exporter
	Auth *configauth.Authentication `mapstructure:"auth"`
	// Headers are the explicit extra headers that should be sent with the exporter
	Headers map[string]string `mapstructure:"headers,omitempty"`
}

// OtlpExporterCfg generates the configuration for a otlp exporter.
func OtlpExporterCfg(endpoint string) *ExporterConfig {
	cfg := ExporterConfig{}
	cfg.Endpoint = endpoint
	return &cfg
}

// OtlpExporterHCPCfg generates the config for an otlp exporter to HCP.
func OtlpExporterHCPCfg(endpoint, resourceID string, authID component.ID) *ExporterConfig {
	// TODO: unfortunately we can't use the exporter config that comes form the otlphttpexporter.Config
	// due to unmarshalling issues. This is unfortunate but for now it's not the end of the world to ship our own config. Leaving this here as a reference
	// to get to the defaultCfg if it's needed.
	//
	// defaultCfg := otlphttpexporter.NewFactory().CreateDefaultConfig().(*otlphttpexporter.Config)
	// defaultCfg.HTTPClientSettings.Endpoint = endpoint
	// defaultCfg.HTTPClientSettings.Auth = &configauth.Authentication{AuthenticatorID: authId}

	cfg := ExporterConfig{
		Headers: make(map[string]string),
	}
	cfg.Endpoint = endpoint
	cfg.Auth = &configauth.Authentication{AuthenticatorID: authID}
	cfg.Headers[channelName] = fmt.Sprintf("%s/%s", channelValue, version.GetHumanVersion())
	cfg.Headers[resourceIDHeader] = resourceID

	return &cfg
}
