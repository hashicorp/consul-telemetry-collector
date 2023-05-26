package envoyreceiver

import (
	"context"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/confignet"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
)

const (
	// ID is the indentifier for the receiver.
	ID = "envoy"

	defaultGRPCEndpoint = "0.0.0.0:9356"
)

// NewFactory creates a new envoy receiver factory.
func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		ID,
		CreateDefaultConfig,
		receiver.WithMetrics(createMetrics, component.StabilityLevelDevelopment),
	)
}

func createMetrics(_ context.Context,
	set receiver.CreateSettings,
	cfg component.Config,
	// nextConsumer is whatever component is next on the pipeline.
	nextConsumer consumer.Metrics) (receiver.Metrics, error) {
	if nextConsumer == nil {
		return nil, component.ErrNilNextConsumer
	}

	envoyCfg := cfg.(*Config)
	envoy := newEnvoyReceiver(set, envoyCfg)

	envoy.registerMetrics(nextConsumer)

	return receiver.Metrics(envoy), nil
}

// CreateDefaultConfig creates the default configuration for receiver.
func CreateDefaultConfig() component.Config {
	return &Config{
		GRPC: &configgrpc.GRPCServerSettings{
			NetAddr: confignet.NetAddr{
				Endpoint:  defaultGRPCEndpoint,
				Transport: "tcp",
			},
			// We almost write 0 bytes, so no need to tune WriteBufferSize.
			ReadBufferSize: 512 * 1024,
			Keepalive: &configgrpc.KeepaliveServerConfig{
				ServerParameters: &configgrpc.KeepaliveServerParameters{
					MaxConnectionIdle: 5 * time.Second,
					MaxConnectionAge:  1 * time.Minute,
					Time:              30 * time.Second,
				},
				EnforcementPolicy: &configgrpc.KeepaliveEnforcementPolicy{
					MinTime: 5 * time.Second,
				},
			},
		},
	}
}
