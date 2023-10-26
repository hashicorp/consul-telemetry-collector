// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package exporters

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
)

func Test_OtlpExporter(t *testing.T) {
	e := &ExporterConfig{
		Endpoint: "foobar",
	}
	conf, err := OtlpExporterCfg(e)
	require.NoError(t, err)

	// Unmarshall and verify
	unmarshalledCfg := &ExporterConfig{}
	err = conf.Unmarshal(unmarshalledCfg)
	require.NoError(t, err)

	require.Equal(t, e.Endpoint, unmarshalledCfg.Endpoint)
	require.Equal(t, e.Compression, defaultCompression)
	require.Equal(t, unmarshalledCfg.Headers["user-agent"], defaultUserAgent)
}

func Test_OtlpExporterHCP(t *testing.T) {
	cfg := OtlpExporterHCPCfg("foobar", "resource-id", component.NewID("foobarid"))
	require.NotNil(t, cfg)

	// Marshall the configuration
	conf := confmap.New()
	err := conf.Marshal(cfg)
	require.NoError(t, err)

	// Unmarshall and verify
	unmarshalledCfg := &ExporterConfig{
		Headers: map[string]string{},
	}
	err = conf.Unmarshal(unmarshalledCfg)
	require.NoError(t, err)

	require.Equal(t, cfg, unmarshalledCfg)
}
