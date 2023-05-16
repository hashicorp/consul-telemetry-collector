package processors

import (
	"testing"

	"github.com/google/uuid"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourceprocessor"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/confmap"
)

func Test_ResourceProcessorCfg(t *testing.T) {
	clusterVal := uuid.NewString()
	cfg := ResourcesProcessorCfg(clusterVal)
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
	require.Equal(t, unmarshalledCfg.AttributesActions[0].Key, "cluster")
	require.Equal(t, unmarshalledCfg.AttributesActions[0].Value, clusterVal)
}
