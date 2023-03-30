package processors

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/processor/batchprocessor"
)

func Test_BatchProcessor(t *testing.T) {
	cfg := BatchProcessorCfg()
	require.NotNil(t, cfg)

	// Marshall the configuration
	conf := confmap.New()
	err := conf.Marshal(cfg)
	require.NoError(t, err)

	//Unmarshall and verify
	unmarshalledCfg := &batchprocessor.Config{}
	conf.Unmarshal(unmarshalledCfg)

	require.Equal(t, cfg, unmarshalledCfg)
}
