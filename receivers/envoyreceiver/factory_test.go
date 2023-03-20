package envoyreceiver

import (
	"context"
	"fmt"
	"testing"

	"github.com/shoenig/test/must"
	"github.com/shoenig/test/portal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/confighttp"
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

	// cfg.GRPC.NetAddr.Endpoint = localEndpoint(t)
	// cfg.HTTP.Endpoint = localEndpoint(t)

	creationSet := receivertest.NewNopCreateSettings()

	mReceiver, err := factory.CreateMetricsReceiver(context.Background(), creationSet, cfg, consumertest.NewNop())
	assert.NotNil(t, mReceiver)
	assert.NoError(t, err)
}

func TestCreateMetricReceiver(t *testing.T) {
	factory := NewFactory()
	defaultGRPCSettings := &configgrpc.GRPCServerSettings{
		NetAddr: confignet.NetAddr{
			Endpoint:  localEndpoint(t),
			Transport: "tcp",
		},
	}
	defaultHTTPSettings := &confighttp.HTTPServerSettings{
		Endpoint: localEndpoint(t),
	}

	_ = defaultHTTPSettings
	_ = defaultGRPCSettings

	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{
			name: "default",
			cfg:  &Config{
				// Protocols: Protocols{
				// 	GRPC: defaultGRPCSettings,
				// 	HTTP: defaultHTTPSettings,
				// },
			},
		},
		// {
		// 	name: "invalid_grpc_address",
		// 	cfg:  &Config{
		// Protocols: Protocols{
		// 	GRPC: &configgrpc.GRPCServerSettings{
		// 		NetAddr: confignet.NetAddr{
		// 			Endpoint:  "327.0.0.1:1122",
		// 			Transport: "tcp",
		// 		},
		// 	},
		// 	HTTP: defaultHTTPSettings,
		// },
		// },
		// wantErr: true,
		// },
		// {
		// 	name: "invalid_http_address",
		// 	cfg:  &Config{
		// Protocols: Protocols{
		// 	GRPC: defaultGRPCSettings,
		// 	HTTP: &confighttp.HTTPServerSettings{
		// 		Endpoint: "327.0.0.1:1122",
		// 	},
		// },
		// },
		// wantErr: true,
		// },
	}
	ctx := context.Background()
	creationSet := receivertest.NewNopCreateSettings()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sink := new(consumertest.MetricsSink)
			mr, err := factory.CreateMetricsReceiver(ctx, creationSet, tt.cfg, sink)
			assert.NoError(t, err)
			require.NotNil(t, mr)
			if tt.wantErr {
				assert.Error(t, mr.Start(context.Background(), componenttest.NewNopHost()))
			} else {
				require.NoError(t, mr.Start(context.Background(), componenttest.NewNopHost()))
				assert.NoError(t, mr.Shutdown(context.Background()))
			}
		})
	}
}

func localEndpoint(t *testing.T) string {
	grabber := portal.New(t, portal.WithAddress("127.0.0.1"))
	port := grabber.One()
	return fmt.Sprintf("127.0.0.1:%d", port)
}
