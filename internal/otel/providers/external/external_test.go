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
)

func Test_InMem(t *testing.T) {
	t.Run("no forwarder", func(t *testing.T) {
		provider := NewProvider(nil)
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
		provider := NewProvider(&config.ExportConfig{
			ID: exporters.BaseOtlpExporterID,
			Exporter: &exporters.ExporterConfig{
				Endpoint: "https://localhost:6060",
			},
		})
		retrieved, err := provider.Retrieve(context.Background(), "", nil)
		test.NoError(t, err)

		conf, err := retrieved.AsConf()
		test.NoError(t, err)
		confMap := conf.ToStringMap()
		exporters := asMap(t, confMap["exporters"])
		otlp := asMap(t, exporters["otlphttp"])
		test.Eq(t, otlp["endpoint"], "https://localhost:6060")
	})
}

func asMap(t *testing.T, a any) map[string]any {
	t.Helper()

	m, ok := a.(map[string]any)
	must.True(t, ok)
	return m
}
