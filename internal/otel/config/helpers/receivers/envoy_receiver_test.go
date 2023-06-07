// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package receivers

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/confmap"

	"github.com/hashicorp/consul-telemetry-collector/receivers/envoyreceiver"
)

func Test_EnvoyReceiver(t *testing.T) {
	cfg := EnvoyReceiverCfg()

	conf := confmap.New()
	err := conf.Marshal(cfg)
	require.NoError(t, err)

	// Unmarshall and verify
	unmarshalledCfg := &envoyreceiver.Config{}
	err = conf.Unmarshal(unmarshalledCfg)
	require.NoError(t, err)

	require.Equal(t, cfg, unmarshalledCfg)
}
