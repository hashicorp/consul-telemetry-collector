// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package processors

import (
	"errors"
	"testing"

	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/resourceprocessor"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/confmap"

	"github.com/hashicorp/consul-telemetry-collector/internal/hcp"
)

func Test_ResourceProcessorCfg(t *testing.T) {
	for name, tc := range map[string]struct {
		mockErr error
	}{
		"Success": {},
		"GetError": {
			mockErr: errors.New("boom"),
		},
	} {
		t.Run(name, func(t *testing.T) {
			clientM := &hcp.MockClient{
				MockMetricAttributes: map[string]string{
					"cluster": "name",
					"org":     "fake-org",
				},
				Err: tc.mockErr,
			}

			cfg := ResourcesProcessorCfg(clientM)
			require.NotNil(t, cfg)

			// Marshal the configuration
			conf := confmap.New()
			err := conf.Marshal(cfg)
			require.NoError(t, err)

			// Unmarshal and verify
			unmarshalledCfg := &resourceprocessor.Config{}
			err = conf.Unmarshal(unmarshalledCfg)
			require.NoError(t, err)

			if tc.mockErr != nil {
				require.Empty(t, unmarshalledCfg.AttributesActions)
				return
			}

			require.NoError(t, unmarshalledCfg.Validate())
			require.Len(t, unmarshalledCfg.AttributesActions, 2)
			for _, action := range unmarshalledCfg.AttributesActions {
				require.Contains(t, []string{"cluster", "org"}, action.Key)
				require.Contains(t, []string{"name", "fake-org"}, action.Value)
			}
		})
	}
}
