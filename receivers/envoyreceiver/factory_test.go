package envoyreceiver

import (
	"context"
	"fmt"
	"testing"

	"github.com/shoenig/test"
	"github.com/shoenig/test/must"
	"github.com/shoenig/test/portal"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/confignet"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/receiver/receivertest"
)

func TestCreateDefaultConfig(t *testing.T) {
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()
	must.NotNil(t, cfg, must.Sprint("failed to create default config"))
	must.NoError(t, componenttest.CheckConfigStruct(cfg))
}

func TestCreateReceiver(t *testing.T) {
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig().(*Config)

	creationSet := receivertest.NewNopCreateSettings()

	mReceiver, err := factory.CreateMetricsReceiver(context.Background(), creationSet, cfg, consumertest.NewNop())
	test.NotNil(t, mReceiver)
	test.NoError(t, err)
}

func TestCreateMetricReceiver(t *testing.T) {
	factory := NewFactory()
	defaultGRPCSettings := &configgrpc.GRPCServerSettings{
		NetAddr: confignet.NetAddr{
			Endpoint:  localEndpoint(t),
			Transport: "tcp",
		},
	}

	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{
			name: "default",
			cfg: &Config{
				GRPC: defaultGRPCSettings,
			},
		},
		{
			name: "invalid_grpc_address",
			cfg: &Config{
				GRPC: &configgrpc.GRPCServerSettings{
					NetAddr: confignet.NetAddr{
						Endpoint:  "327.0.0.1:1122",
						Transport: "tcp",
					},
				},
			},
			wantErr: true,
		},
	}

	ctx := context.Background()
	creationSet := receivertest.NewNopCreateSettings()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sink := new(consumertest.MetricsSink)
			mr, err := factory.CreateMetricsReceiver(ctx, creationSet, tt.cfg, sink)
			test.NoError(t, err)
			must.NotNil(t, mr)
			startErr := mr.Start(context.Background(), componenttest.NewNopHost())
			if tt.wantErr {
				test.Error(t, startErr)
			} else {
				must.NoError(t, startErr)
				test.NoError(t, mr.Shutdown(context.Background()))
			}
		})
	}
}

func localEndpoint(t *testing.T) string {
	t.Helper()

	grabber := portal.New(t, portal.WithAddress("127.0.0.1"))
	port := grabber.One()
	return fmt.Sprintf("127.0.0.1:%d", port)
}
