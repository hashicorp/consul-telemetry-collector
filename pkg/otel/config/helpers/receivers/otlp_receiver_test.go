package receivers

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/confmap"
)

func Test_OtlpReceiver(t *testing.T) {
	id, cfg := OtlpReceiverCfg()

	require.NotEmpty(t, id)
	require.NotNil(t, cfg)
	conf := confmap.New()
	err := conf.Marshal(cfg)

	require.NoError(t, err)
	retrieved, _ := confmap.NewRetrieved(conf.ToStringMap())
	require.NotNil(t, retrieved)
}