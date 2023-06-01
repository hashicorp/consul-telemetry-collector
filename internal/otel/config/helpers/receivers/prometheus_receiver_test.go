package receivers

import (
	"fmt"
	"testing"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/confmap"
)

func Test_PrometheusReceiverCfg(t *testing.T) {
	cfg := PrometheusReceiverCfg()

	conf := confmap.New()
	err := conf.Marshal(cfg)
	require.NoError(t, err)

	// Unmarshall and verify
	unmarshalledCfg := &prometheusreceiver.Config{}
	err = unmarshalledCfg.Unmarshal(conf)
	require.NoError(t, err)

	fmt.Println(cfg)

	require.NoError(t, unmarshalledCfg.Validate())
	require.NotNil(t, unmarshalledCfg.PrometheusConfig)
	require.Len(t, unmarshalledCfg.PrometheusConfig.ScrapeConfigs, 1)
}
