package processors

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/confmap"

	"github.com/hashicorp/consul-telemetry-collector/internal/hcp"
)

func Test_FilterProcessor(t *testing.T) {
	clientM := &hcp.MockClient{
		MockMetricFilters: []string{
			"^consul.consul.envoy.connection$",
			"^envoy.*connection$",
			"[a-z",
		},
	}
	cfg := FilterProcessorCfg(clientM)
	require.NotNil(t, cfg)

	// Marshall the configuration
	conf := confmap.New()
	err := conf.Marshal(cfg)
	require.NoError(t, err)

	// Unmarshall and verify
	unmarshalledCfg := &FilterProcessorConfig{}
	err = conf.Unmarshal(&unmarshalledCfg)
	require.NoError(t, err)
	require.Equal(t, []string{"^consul.consul.envoy.connection$", "^envoy.*connection$"},
		unmarshalledCfg.Metrics.Include.MetricNames)

	require.Equal(t, cfg, unmarshalledCfg)
}
