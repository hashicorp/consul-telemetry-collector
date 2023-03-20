package envoyreceiver

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/shoenig/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/confignet"
	"go.opentelemetry.io/collector/confmap/confmaptest"
)

func TestUnmarshalDefaultConfig(t *testing.T) {
	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "default.yaml"))
	require.NoError(t, err)
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()
	assert.NoError(t, component.UnmarshalConfig(cm, cfg))
	assert.Equal(t, factory.CreateDefaultConfig(), cfg)
}

func TestUnmarshalConfigUnix(t *testing.T) {
	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "uds.yaml"))
	require.NoError(t, err)
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()
	test.NoError(t, component.UnmarshalConfig(cm, cfg))

	marshalCfg := cfg.(*Config)

	test.Eq(t,
		&Config{
			GRPC: &configgrpc.GRPCServerSettings{
				NetAddr: confignet.NetAddr{
					Endpoint:  "/tmp/grpc_otlp.sock",
					Transport: "unix",
				},
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
		}, marshalCfg)
}
