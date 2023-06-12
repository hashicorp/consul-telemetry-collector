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

func TestBuilder_Gauge(t *testing.T) {
	goldenBytes, err := os.ReadFile("testdata/gauge.golden")
	must.NoError(t, err)
	goldenGauges, err := new(pmetric.JSONUnmarshaler).UnmarshalMetrics(goldenBytes)
	must.NoError(t, err)
	f := "testdata/gauge"
	bytes, err := os.ReadFile(f)
	must.NoError(t, err)

	labels := map[string]string{
		"name":    uuid.NewString(),
		"cluster": uuid.NewString(),
	}

	promGauges := make([]*prompb.MetricFamily, 0)
	must.NoError(t, json.Unmarshal(bytes, &promGauges))

	b := NewBuilder(labels)
	for _, gauge := range promGauges {
		b.AddGauge(gauge)
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
	must.Length(t, 3, metricSlice)

	for i := 0; i < metricSlice.Len(); i++ {
		metric := metricSlice.At(0)
		must.Eq(t, pmetric.MetricTypeGauge, metric.Type(), must.Sprintf("metric types don't match %s!=%s",
			pmetric.MetricTypeGauge.String(), metric.Type().String()))
	}

	goldenGaugeSlice := goldenGauges.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics()

	goldenTestGauge := flattenGauge(goldenGaugeSlice)
	gauges := flattenGauge(metricSlice)

	for _, ct := range gauges {
		must.Contains[testmetric](t, ct, ContainsTestMetric(goldenTestGauge), must.Sprintf("metric %s not found",
			ct.name))
	}
}

func flattenGauge(ms pmetric.MetricSlice) []testmetric {
	gauges := make([]testmetric, 0)
	for i := 0; i < ms.Len(); i++ {
		metric := ms.At(i)
		for j := 0; j < metric.Gauge().DataPoints().Len(); j++ {
			dp := metric.Gauge().DataPoints().At(j)
			attrs := map[string]string{}
			dp.Attributes().Range(func(k string, v pcommon.Value) bool {
				k = strings.ReplaceAll(k, ".", "_")
				attrs[k] = v.Str()
				return true
			})
			gauges = append(gauges, testmetric{
				name:       strings.ReplaceAll(metric.Name(), ".", "_"),
				val:        dp.DoubleValue(),
				attributes: attrs,
			})
		}
	}
	return gauges
}
