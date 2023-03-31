package exporters

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
)

func Test_OtlpExporter(t *testing.T) {
	cfg := OtlpExporterCfg("foobar")
	require.NotNil(t, cfg)

	// Marshall the configuration
	conf := confmap.New()
	err := conf.Marshal(cfg)
	require.NoError(t, err)

	// Unmarshall and verify
	unmarshalledCfg := &ExporterConfig{}
	err = conf.Unmarshal(unmarshalledCfg)
	require.NoError(t, err)

	require.Equal(t, cfg, unmarshalledCfg)
}

func Test_OtlpExporterHCP(t *testing.T) {
	cfg := OtlpExporterHCPCfg("foobar", component.NewID("foobarid"))
	require.NotNil(t, cfg)

	// Marshall the configuration
	conf := confmap.New()
	err := conf.Marshal(cfg)
	require.NoError(t, err)

	// Unmarshall and verify
	unmarshalledCfg := &ExporterConfig{}
	err = conf.Unmarshal(unmarshalledCfg)
	require.NoError(t, err)
	require.Equal(t, cfg, unmarshalledCfg)
}
