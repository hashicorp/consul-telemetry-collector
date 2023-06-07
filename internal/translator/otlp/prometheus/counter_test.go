// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package prometheus

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	prompb "github.com/prometheus/client_model/go"
	"github.com/shoenig/test/must"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

type testmetric struct {
	name       string
	val        float64
	attributes map[string]string
}

func TestBuilder_Counter(t *testing.T) {
	goldenBytes, err := os.ReadFile("testdata/counter.golden")
	must.NoError(t, err)
	goldenMetrics, err := new(pmetric.JSONUnmarshaler).UnmarshalMetrics(goldenBytes)
	must.NoError(t, err)
	f := "testdata/counter"
	labels := map[string]string{
		"name":    uuid.NewString(),
		"cluster": uuid.NewString(),
	}

	counterBytes, err := os.ReadFile(f)
	must.NoError(t, err)

	promCounter := make([]*prompb.MetricFamily, 0)
	must.NoError(t, json.Unmarshal(counterBytes, &promCounter))

	b := NewBuilder(labels)
	for _, counter := range promCounter {
		b.AddCounter(counter)
	}

	md := b.Build()

	must.Length(t, 1, md.ResourceMetrics())
	for k, v := range labels {
		val, ok := md.ResourceMetrics().At(0).Resource().Attributes().Get(k)
		must.True(t, ok)
		must.Eq(t, v, val.AsString())
	}
	md.ResourceMetrics().At(0).Resource().Attributes().Range(func(k string, v pcommon.Value) bool {
		val, ok := labels[k]
		must.True(t, ok)
		must.Eq(t, v.AsString(), val)
		return true
	})

	must.Length(t, 1, md.ResourceMetrics().At(0).ScopeMetrics())
	metricSlice := md.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics()
	must.Length(t, 4, metricSlice)
	must.Eq(t, 4, goldenMetrics.DataPointCount())
	goldenMetricSlice := goldenMetrics.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics()

	for i := 0; i < metricSlice.Len(); i++ {
		metric := metricSlice.At(0)
		must.Eq(t, pmetric.MetricTypeSum, metric.Type())
		must.True(t, metric.Sum().IsMonotonic())
	}

	goldenCounters := flattenCounter(goldenMetricSlice)
	counters := flattenCounter(metricSlice)

	for _, ct := range counters {
		must.Contains[testmetric](t, ct, ContainsTestMetric(goldenCounters), must.Sprintf("does not contain %v",
			ct))
	}
}

type ContainsFunc[T any] func(T) bool

func (c ContainsFunc[T]) Contains(v T) bool {
	return (c)(v)
}

func ContainsTestMetric(metric []testmetric) ContainsFunc[testmetric] {
	return func(t testmetric) bool {
		for _, counter := range metric {
			if counter.name == t.name && counter.val == t.val && len(counter.attributes) == len(t.attributes) {
				for k, v := range counter.attributes {
					val, ok := t.attributes[k]
					if !ok {
						return false
					}
					if val != v {
						return false
					}
				}
				return true
			}
		}
		return false
	}
}

func flattenCounter(ms pmetric.MetricSlice) []testmetric {
	counters := make([]testmetric, 0)
	for i := 0; i < ms.Len(); i++ {
		metric := ms.At(i)
		for j := 0; j < metric.Sum().DataPoints().Len(); j++ {
			dp := metric.Sum().DataPoints().At(j)
			attrs := map[string]string{}
			dp.Attributes().Range(func(k string, v pcommon.Value) bool {
				k = strings.ReplaceAll(k, ".", "_")
				attrs[k] = v.Str()
				return true
			})
			counters = append(counters, testmetric{
				name:       strings.ReplaceAll(metric.Name(), ".", "_"),
				val:        dp.DoubleValue(),
				attributes: attrs,
			})
		}
	}
	return counters
}
