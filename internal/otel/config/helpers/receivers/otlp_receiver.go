// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package receivers holds the type of receivers that consul telemetery supports
package receivers

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configauth"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"

	"github.com/hashicorp/consul-telemetry-collector/internal/otel/config/helpers/types"
)

const otlpReceiverName = "otlp"

// OtlpReceiverID is the component id of the otlp receiver.
var OtlpReceiverID component.ID = component.NewID(otlpReceiverName)

// Protocols is the configuration for the supported protocols.
type Protocols struct {
	GRPC *configgrpc.GRPCServerSettings `mapstructure:"grpc,omitempty"`
	HTTP *HTTPConfig                    `mapstructure:"http,omitempty"`
}

// HTTPConfig is the HTTPConfig type used to build the http settings for the receiver
type HTTPConfig struct {
	// Endpoint configures the listening address for the server.
	Endpoint string `mapstructure:"endpoint"`

	// TLSSetting struct exposes TLS client configuration.
	TLSSetting *types.TLSServerSetting `mapstructure:"tls"`

	// CORS configures the server for HTTP cross-origin resource sharing (CORS).
	CORS *confighttp.CORSSettings `mapstructure:"cors"`

	// Auth for this receiver
	Auth *configauth.Authentication `mapstructure:"auth"`

	// MaxRequestBodySize sets the maximum request body size in bytes
	MaxRequestBodySize int64 `mapstructure:"max_request_body_size"`

	// IncludeMetadata propagates the client metadata from the incoming requests to the downstream consumers
	// Experimental: *NOTE* this option is subject to change or removal in the future.
	IncludeMetadata bool `mapstructure:"include_metadata"`

	// Additional headers attached to each HTTP response sent to the client.
	// Header values are opaque since they may be sensitive.
	ResponseHeaders map[string]string `mapstructure:"response_headers"`
	// The URL path to receive traces on. If omitted "/v1/traces" will be used.
	TracesURLPath string `mapstructure:"traces_url_path,omitempty"`

	// The URL path to receive metrics on. If omitted "/v1/metrics" will be used.
	MetricsURLPath string `mapstructure:"metrics_url_path,omitempty"`

	// The URL path to receive logs on. If omitted "/v1/logs" will be used.
	LogsURLPath string `mapstructure:"logs_url_path,omitempty"`
}

// OtlpReceiverConfig defines configuration for OTLP receiver.
type OtlpReceiverConfig struct {
	// Protocols is the configuration for the supported protocols, currently gRPC and HTTP (Proto and JSON).
	Protocols `mapstructure:"protocols"`
}

// OtlpReceiverCfg  generates the config for an otlp receiver.
func OtlpReceiverCfg() *OtlpReceiverConfig {
	defaults := otlpreceiver.NewFactory().CreateDefaultConfig().(*otlpreceiver.Config)

	httpConfig := HTTPConfig{
		Endpoint: defaults.HTTP.Endpoint,

		ResponseHeaders: make(map[string]string),
		TracesURLPath:   defaults.HTTP.TracesURLPath,
		LogsURLPath:     defaults.HTTP.LogsURLPath,
		MetricsURLPath:  defaults.HTTP.MetricsURLPath,
	}

	return &OtlpReceiverConfig{
		Protocols: Protocols{
			HTTP: &httpConfig,
		},
	}
}
