package processors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/confmap"

	"github.com/hashicorp/consul-telemetry-collector/internal/hcp"
)

func Test_FilterProcessor(t *testing.T) {
	testcases := map[string]struct {
		mockErr error
	}{
		"Success": {},
		"FilterGetError": {
			mockErr: errors.New("boom"),
		},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			clientM := &hcp.MockClient{
				MockMetricFilters: []string{
					"^consul.consul.envoy.connection$",
					"^envoy.*connection$",
					"[a-z",
				},
				Err: tc.mockErr,
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
			if tc.mockErr == nil {
				require.Equal(t, []string{"^consul.consul.envoy.connection$", "^envoy.*connection$"},
					unmarshalledCfg.Metrics.Include.MetricNames)
			}

			require.Equal(t, cfg, unmarshalledCfg)
		})
	}

}
