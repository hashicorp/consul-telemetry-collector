package exporters

import (
	"fmt"
	"os"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configauth"
	"go.opentelemetry.io/collector/config/configtls"

	"github.com/hashicorp/consul-telemetry-collector/internal/version"
)

const (
	// otlpHTTPExporterName is the component.ID value used by the otlphttp exporter.
	otlpHTTPExporterName = "otlphttp"
	channelName          = "x-channel"
	channelValue         = "consul-telemetry-collector"
	resourceIDHeader     = "x-hcp-resource-id"

	userAgentHeader  = "user-agent"
	defaultUserAgent = "Go-http-client/1.1"

	envVarOtlpExporterTLS = "OTLP_EXPORTER_TLS"
	tlsSettingInsecure    = "insecure"
	tlsSettingDisabled    = "disabled"
)

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

	// TLSSetting struct exposes TLS client configuration.
	TLSSetting configtls.TLSClientSetting `mapstructure:"tls"`

	// The compression key for supported compression types within collector.
	Compression string `mapstructure:"compression"`
}

// OtlpExporterCfg generates the configuration for a otlp exporter.
func OtlpExporterCfg(endpoint string) *ExporterConfig {
	cfg := ExporterConfig{
		Compression: "none",
		Headers: map[string]string{
			userAgentHeader: defaultUserAgent,
		},
	}
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
		Headers: map[string]string{
			channelName:      fmt.Sprintf("%s/%s", channelValue, version.GetHumanVersion()),
			resourceIDHeader: resourceID,
			userAgentHeader:  defaultUserAgent,
		},
		Auth:        &configauth.Authentication{AuthenticatorID: authID},
		Endpoint:    endpoint,
		TLSSetting:  tlsConfigForSetting(),
		Compression: "none",
	}

	return &cfg
}

func tlsConfigForSetting() configtls.TLSClientSetting {
	setting := os.Getenv(envVarOtlpExporterTLS)
	switch setting {
	// case tlsSettingDisabled:
	// 	return configtls.TLSClientSetting{Insecure: true}
	case tlsSettingInsecure:
		return configtls.TLSClientSetting{InsecureSkipVerify: true}
	default:
		return configtls.TLSClientSetting{}
	}
}
