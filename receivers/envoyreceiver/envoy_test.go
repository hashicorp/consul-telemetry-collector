package envoyreceiver

import (
	"path/filepath"
	"testing"

	"github.com/shoenig/test"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confignet"
	"go.opentelemetry.io/collector/confmap/confmaptest"
)

func TestUnmarshalDefaultConfig(t *testing.T) {
	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "default.yaml"))
	require.NoError(t, err)
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()
	test.NoError(t, component.UnmarshalConfig(cm, cfg))
	test.Eq(t, factory.CreateDefaultConfig(), cfg)
}

func TestUnmarshalConfigUnix(t *testing.T) {
	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "uds.yaml"))
	require.NoError(t, err)
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()
	test.NoError(t, component.UnmarshalConfig(cm, cfg))

	marshalCfg := cfg.(*Config)

	actualComponentCfg := factory.CreateDefaultConfig()
	actualCfg := actualComponentCfg.(*Config)
	actualCfg.GRPC.NetAddr = confignet.NetAddr{
		Endpoint:  "/tmp/grpc_otlp.sock",
		Transport: "unix",
	}
	test.Eq(t, actualCfg, marshalCfg)
}
