// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package external

import (
	"context"
	"testing"

	"github.com/shoenig/test"
	"github.com/shoenig/test/must"
)

func Test_InMem(t *testing.T) {
	provider := NewProvider("https://localhost:6060", "")
	retrieved, err := provider.Retrieve(context.Background(), "", nil)
	test.NoError(t, err)

	conf, err := retrieved.AsConf()
	test.NoError(t, err)
	confMap := conf.ToStringMap()
	exporters := asMap(t, confMap["exporters"])
	otlp := asMap(t, exporters["otlphttp"])
	test.Eq(t, otlp["endpoint"], "https://localhost:6060")
}

func Test_InMemWithOverdies(t *testing.T) {
	provider := NewProvider("https://localhost:6060", "./testdata/test.yaml")
	overrides, err := loadOveride("./testdata/test.yaml")
	test.NoError(t, err)
	retrieved, err := provider.Retrieve(context.Background(), "", nil)
	test.NoError(t, err)

	conf, err := retrieved.AsConf()
	test.NoError(t, err)
	confMap := conf.ToStringMap()
	overridesMap := overrides.ToStringMap()

	processors := asMap(t, confMap["processors"])
	overridesProcessors := asMap(t, overridesMap["processors"])

	actualLimiter := asMap(t, processors["memory_limiter"])
	expectedLimiter := asMap(t, overridesProcessors["memory_limiter"])

	test.Eq(t, actualLimiter["limit_percentage"], expectedLimiter["limit_percentage"])
	test.Eq(t, actualLimiter["spike_limit_percentage"], expectedLimiter["spike_limit_percentage"])

	test.MapContainsKey(t, actualLimiter, "spike_limit_mib")
	test.MapContainsKey(t, actualLimiter, "limit_mib")

	test.MapNotContainsKey(t, expectedLimiter, "spike_limit_mib")
	test.MapNotContainsKey(t, expectedLimiter, "limit_mib")

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
