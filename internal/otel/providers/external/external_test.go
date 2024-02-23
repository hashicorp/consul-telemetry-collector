// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package external

import (
	"context"
	"testing"

	"github.com/shoenig/test"
	"github.com/shoenig/test/must"

	"github.com/hashicorp/consul-telemetry-collector/internal/otel/config"
	"github.com/hashicorp/consul-telemetry-collector/internal/otel/config/helpers/exporters"
	"github.com/hashicorp/consul-telemetry-collector/internal/otel/providers"
)

func Test_InMem(t *testing.T) {
	t.Run("no forwarder", func(t *testing.T) {
		provider := NewProvider(nil, providers.SharedParams{})
		retrieved, err := provider.Retrieve(context.Background(), "", nil)
		test.NoError(t, err)

		conf, err := retrieved.AsConf()
		test.NoError(t, err)
		confMap := conf.ToStringMap()
		exporters := asMap(t, confMap["exporters"])
		otlp, ok := exporters["otlphttp"]
		test.False(t, ok)
		test.Nil(t, otlp)
	})

	t.Run("with forwarder", func(t *testing.T) {
		provider := NewProvider(&config.ExporterConfig{
			ID: exporters.BaseOtlpExporterID,
			Exporter: &exporters.ExporterConfig{
				Endpoint: "https://localhost:6060",
			},
		}, providers.SharedParams{})
		retrieved, err := provider.Retrieve(context.Background(), "", nil)
		test.NoError(t, err)

		conf, err := retrieved.AsConf()
		test.NoError(t, err)
		confMap := conf.ToStringMap()
		exporters := asMap(t, confMap["exporters"])
		otlp := asMap(t, exporters["otlphttp"])
		test.Eq(t, otlp["endpoint"], "https://localhost:6060")
	})

	t.Run("with tls grpc forwarder", func(t *testing.T) {
		provider := NewProvider(&config.ExporterConfig{
			ID: exporters.GRPCOtlpExporterID,
			Exporter: &exporters.ExporterConfig{
				Endpoint: "https://localhost:6060",
			},
		}, providers.SharedParams{})
		retrieved, err := provider.Retrieve(context.Background(), "", nil)
		test.NoError(t, err)

		conf, err := retrieved.AsConf()
		test.NoError(t, err)
		confMap := conf.ToStringMap()
		exporters := asMap(t, confMap["exporters"])
		otlp := asMap(t, exporters["otlp"])
		test.Eq(t, otlp["endpoint"], "https://localhost:6060")
		test.Eq(t, otlp["tls"], nil)
	})

	t.Run("with non-tls grpc forwarder", func(t *testing.T) {
		provider := NewProvider(&config.ExporterConfig{
			ID: exporters.GRPCOtlpExporterID,
			Exporter: &exporters.ExporterConfig{
				Endpoint: "http://localhost:6060",
			},
		}, providers.SharedParams{})
		retrieved, err := provider.Retrieve(context.Background(), "", nil)
		test.NoError(t, err)

		conf, err := retrieved.AsConf()
		test.NoError(t, err)
		confMap := conf.ToStringMap()
		exporters := asMap(t, confMap["exporters"])
		otlp := asMap(t, exporters["otlp"])
		test.Eq(t, otlp["endpoint"], "http://localhost:6060")
		tlsSetting := asMap(t, otlp["tls"])
		test.Eq(t, tlsSetting["insecure"], true)
		test.Eq(t, tlsSetting["insecure_skip_verify"], false)
	})
}

func asMap(t *testing.T, a any) map[string]any {
	t.Helper()

	m, ok := a.(map[string]any)
	must.True(t, ok)
	return m
}
