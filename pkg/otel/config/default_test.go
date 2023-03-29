package config

import (
	"testing"

	"github.com/hashicorp/consul-telemetry-collector/pkg/otel/config/helpers/exporters"
	"github.com/hashicorp/consul-telemetry-collector/pkg/otel/config/helpers/extensions"
	"github.com/hashicorp/consul-telemetry-collector/pkg/otel/config/helpers/processors"
	"github.com/hashicorp/consul-telemetry-collector/pkg/otel/config/helpers/receivers"
	"github.com/stretchr/testify/require"
)

const (
	otlpExporterEP = "foobar.ca"
	hcpExporterEP  = "hcp.dev.to"
)

func baseAssertions(t *testing.T, cfg *Config, iCfg *IntermediateConfig) {
	require.Empty(t, cfg.Service.Extensions)
	require.Empty(t, cfg.Service.Pipelines)
	logID, logCfg := exporters.LogExporterCfg()
	val, ok := cfg.Exporters[logID]
	require.Equal(t, ok, true)
	require.Equal(t, val, logCfg)

	otlpReceiverID, otlpReceiverCfg := receivers.OtlpReceiverCfg()
	val, ok = cfg.Receivers[otlpReceiverID]
	require.Equal(t, ok, true)
	require.Equal(t, val, otlpReceiverCfg)

	batchID, batchCfg := processors.BatchProcessorCfg()
	val, ok = cfg.Processors[batchID]
	require.Equal(t, ok, true)
	require.Equal(t, val, batchCfg)

	memLimitID, memLimitCfg := processors.MemoryLimiterCfg()
	val, ok = cfg.Processors[memLimitID]
	require.Equal(t, ok, true)
	require.Equal(t, val, memLimitCfg)

	ballastID, ballastCfg := extensions.BallastCfg()
	val, ok = cfg.Extensions[ballastID]
	require.Equal(t, ok, true)
	require.Equal(t, val, ballastCfg)

}

type mockTelemteryClient struct {
	endpoint    string
	expectedErr error
	filter      []string
}

func (m mockTelemteryClient) MetricsEndpoint() (string, error) {
	return m.endpoint, m.expectedErr
}

func (m mockTelemteryClient) MetricFilters() ([]string, error) {
	return m.filter, m.expectedErr
}

func Test_DefaultConfig(t *testing.T) {
	for name, tc := range map[string]struct {
		params     *DefaultParams
		assertions func(*testing.T, *Config, *IntermediateConfig)
	}{
		"base": {
			params: &DefaultParams{},
			assertions: func(t *testing.T, cfg *Config, iCfg *IntermediateConfig) {
				baseAssertions(t, cfg, iCfg)
				require.Equal(t, len(iCfg.Exporters), 1)
				require.Equal(t, len(iCfg.Receivers), 1)
				require.Equal(t, len(iCfg.Processors), 2)
				require.Equal(t, len(iCfg.Extensions), 1)

			},
		},
		"IncludedOtel": {
			params: &DefaultParams{
				OtlpHTTPEndpoint: otlpExporterEP,
			},
			assertions: func(t *testing.T, cfg *Config, iCfg *IntermediateConfig) {
				baseAssertions(t, cfg, iCfg)
				require.Equal(t, len(iCfg.Exporters), 2)
				baseExporterID, baseExporterCfg := exporters.OtlpExporterCfg(otlpExporterEP)
				val, ok := cfg.Exporters[baseExporterID]
				require.Equal(t, ok, true)
				require.Equal(t, val, baseExporterCfg)
			},
		},
		"IncludedOtelHTP": {
			params: &DefaultParams{
				OtlpHTTPEndpoint: otlpExporterEP,
				Client:           mockTelemteryClient{endpoint: hcpExporterEP},
				ClientID:         "id",
				ClientSecret:     "sec",
			},
			assertions: func(t *testing.T, cfg *Config, iCfg *IntermediateConfig) {
				baseAssertions(t, cfg, iCfg)
				require.Equal(t, len(iCfg.Exporters), 3)
				baseExporterID, baseExporterCfg := exporters.OtlpExporterCfg(otlpExporterEP)
				val, ok := cfg.Exporters[baseExporterID]
				require.Equal(t, ok, true)
				require.Equal(t, val, baseExporterCfg)

				hcpExporterID, hcpExporterCfg := exporters.OtlpExporterHCPCfg(hcpExporterEP, extensions.Oauth2ClientID)
				val, ok = cfg.Exporters[hcpExporterID]
				require.Equal(t, ok, true)
				require.Equal(t, val, hcpExporterCfg)
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			c, intermediateCfg, err := DefaultConfig(tc.params)
			require.NoError(t, err)
			tc.assertions(t, c, intermediateCfg)
		})
	}

}
