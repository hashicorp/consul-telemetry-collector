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
	typeStr = "envoy"

	defaultGRPCEndpoint = "0.0.0.0:9356"
)

func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		typeStr,
		createDefaultConfig,
		receiver.WithMetrics(func(ctx context.Context, settings receiver.CreateSettings, config component.Config, metrics consumer.Metrics) (receiver.Metrics, error) {
			return receiver.Metrics(&envoyReceiver{}), nil
		}, component.StabilityLevelDevelopment),
	)
}

// createDefaultConfig creates the default configuration for receiver.
func createDefaultConfig() component.Config {
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
