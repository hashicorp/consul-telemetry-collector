package extensions

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/extension/ballastextension"
)

func Test_BallastExtension(t *testing.T) {
	cfg := BallastCfg()
	require.NotNil(t, cfg)

	// Marshall the configuration
	conf := confmap.New()
	err := conf.Marshal(cfg)
	require.NoError(t, err)

	//Unmarshall and verify
	unmarshalledCfg := &ballastextension.Config{}
	conf.Unmarshal(unmarshalledCfg)

	require.Equal(t, cfg, unmarshalledCfg)
}
