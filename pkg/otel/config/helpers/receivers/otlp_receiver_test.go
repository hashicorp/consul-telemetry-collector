package receivers

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/confmap"
)

func Test_OtlpReceiver(t *testing.T) {
	cfg := OtlpReceiverCfg()

	conf := confmap.New()
	err := conf.Marshal(cfg)
	require.NoError(t, err)

	// Unmarshall and verify
	unmarshalledCfg := &OtlpReceiverConfig{}
	err = conf.Unmarshal(unmarshalledCfg)
	require.NoError(t, err)

	require.Equal(t, cfg, unmarshalledCfg)
}
