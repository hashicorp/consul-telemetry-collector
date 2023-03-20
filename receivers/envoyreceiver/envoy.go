package envoyreceiver

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/receiver"
)

type envoyReceiver struct {
	component.Component
}

var _ receiver.Metrics = (*envoyReceiver)(nil)

type Config struct {
	GRPC *configgrpc.GRPCServerSettings `mapstructure:"grpc"`
}

func (r *envoyReceiver) Start(_ context.Context, host component.Host) error {
	return nil
}

func (r *envoyReceiver) Shutdown(ctx context.Context) error {
	return nil
}
