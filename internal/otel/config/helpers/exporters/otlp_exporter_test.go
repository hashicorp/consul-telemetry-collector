// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package exporters

import (
	"testing"

	"github.com/shoenig/test/must"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/exporter/otlpexporter"
)

func Test_OtlpHTTPExporter(t *testing.T) {
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

func Test_OtlpHTTPExporterHCP(t *testing.T) {
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
	require.Nil(t, cfg.TLSSetting)

	require.Equal(t, cfg, unmarshalledCfg)
}

func Test_OtlpExporter(t *testing.T) {
	tests := map[string]struct {
		cfg *ExporterConfig
		env func(t *testing.T)
		tls bool
	}{
		"default": {
			cfg: &ExporterConfig{
				Endpoint: "http://foobar",
			},
		},
		"headers": {
			cfg: &ExporterConfig{
				Endpoint: "http://foobar",
				Headers: map[string]string{
					"a": "b",
				},
			},
		},
		"timeout": {
			cfg: &ExporterConfig{
				Endpoint: "http://foobar",
			},
		},
		"tls": {
			cfg: &ExporterConfig{
				Endpoint: "https://foobar",
			},
			tls: true,
		},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			if testcase.env != nil {
				testcase.env(t)
			}
			conf, err := OtlpExporterCfg(testcase.cfg)
			must.NoError(t, err)

			otlpCfg := &otlpexporter.Config{}
			must.NoError(t, conf.Unmarshal(otlpCfg))
			must.NoError(t, otlpCfg.Validate())

			must.Eq(t, testcase.cfg.Compression, string(otlpCfg.Compression))
			must.Eq(t, testcase.cfg.Endpoint, otlpCfg.Endpoint)
			must.Eq(t, testcase.cfg.Auth, otlpCfg.Auth)

			// when creating a grpc conn the collector calls this to turn the tls settings into a *tls.Config. We're expecting this to be nil on a plaintext endpoint so that we properly do not do TLS
			tlsCfg, err := otlpCfg.TLSSetting.LoadTLSConfig()
			must.NoError(t, err)
			if testcase.tls {
				must.NotNil(t, tlsCfg)
			} else {
				must.True(t, testcase.cfg.TLSSetting.Insecure)
				must.Nil(t, tlsCfg)
			}

			must.MapContainsKey(t, testcase.cfg.Headers, userAgentHeader)
		})
	}
}
