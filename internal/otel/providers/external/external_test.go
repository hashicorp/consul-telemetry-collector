// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package external

import (
	"context"
	"testing"

	"github.com/shoenig/test"
	"github.com/shoenig/test/must"

	"github.com/hashicorp/consul-telemetry-collector/internal/otel/config/helpers/exporters"
)

func Test_InMem(t *testing.T) {
	provider := NewProvider(exporters.BaseOtlpExporterID, &exporters.ExporterConfig{
		Endpoint: "https://localhost:6060",
	})
	retrieved, err := provider.Retrieve(context.Background(), "", nil)
	test.NoError(t, err)

	conf, err := retrieved.AsConf()
	test.NoError(t, err)
	confMap := conf.ToStringMap()
	exporters := asMap(t, confMap["exporters"])
	otlp := asMap(t, exporters["otlphttp"])
	test.Eq(t, otlp["endpoint"], "https://localhost:6060")
}

func asMap(t *testing.T, a any) map[string]any {
	t.Helper()

	m, ok := a.(map[string]any)
	must.True(t, ok)
	return m
}
