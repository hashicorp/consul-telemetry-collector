package processors

import (
	"testing"

	"github.com/google/uuid"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourceprocessor"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/confmap"
)

func Test_ResourceProcessorCfg(t *testing.T) {
	cfg := ResourcesProcessorCfg(uuid.NewString())
	require.NotNil(t, cfg)

	// Marshall the configuration
	conf := confmap.New()
	err := conf.Marshal(cfg)
	require.NoError(t, err)

	// Unmarshall and verify
	unmarshalledCfg := &resourceprocessor.Config{}
	err = conf.Unmarshal(unmarshalledCfg)
	require.NoError(t, err)

	require.NoError(t, unmarshalledCfg.Validate())
	require.Len(t, unmarshalledCfg.AttributesActions, 1)

	mshCfg := confmap.New()
	require.NoError(t, mshCfg.Marshal(unmarshalledCfg))
	// make sure that the initially marshalled cfg matches the marshal'd resourceprocessor match
	require.Equal(t, mshCfg.AllKeys(), conf.AllKeys())
	require.Equal(t, mshCfg.Get("action"), conf.Get("action"))
	require.Equal(t, mshCfg.Get("key"), conf.Get("key"))
	require.Equal(t, mshCfg.Get("value"), conf.Get("value"))
}
