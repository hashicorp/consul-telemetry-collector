package exporters

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/confmap"
)

func Test_LoggingExporter(t *testing.T) {
	cfg := LogExporterCfg()
	require.NotNil(t, cfg)

	// Marshall the configuration
	conf := confmap.New()
	err := conf.Marshal(cfg)
	require.NoError(t, err)

	//Unmarshall and verify
	unmarshalledCfg := &LoggingConfig{}
	conf.Unmarshal(unmarshalledCfg)

	require.Equal(t, cfg, unmarshalledCfg)
}
