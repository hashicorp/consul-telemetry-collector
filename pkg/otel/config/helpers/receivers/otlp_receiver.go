// Package receivers holds the type of receivers that consul telemetery supports
package receivers

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"

	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/confighttp"
)

const otlpReceiverName = "otlp"

// OtlpReceiverID is the component id of the otlp receiver
var OtlpReceiverID component.ID = component.NewID(otlpReceiverName)

// Protocols is the configuration for the supported protocols.
type Protocols struct {
	GRPC *configgrpc.GRPCServerSettings `mapstructure:"grpc,omitempty"`
	HTTP *confighttp.HTTPServerSettings `mapstructure:"http,omitempty"`
}

// OtlpReceiverConfig defines configuration for OTLP receiver.
type OtlpReceiverConfig struct {
	// Protocols is the configuration for the supported protocols, currently gRPC and HTTP (Proto and JSON).
	Protocols `mapstructure:"protocols"`
}

// OtlpReceiverCfg  generates the config for an otlp receiver
func OtlpReceiverCfg() (component.ID, *OtlpReceiverConfig) {
	// cfg := otlpreceiver.Config{}
	// cfg.HTTP = &confighttp.HTTPServerSettings{}
	defaults := otlpreceiver.NewFactory().CreateDefaultConfig().(*otlpreceiver.Config)

	return OtlpReceiverID, &OtlpReceiverConfig{
		Protocols: Protocols{
			HTTP: defaults.HTTP,
		},
	}
}
