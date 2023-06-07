// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package prometheus

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/google/uuid"
	prompb "github.com/prometheus/client_model/go"
	"github.com/shoenig/test/must"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

func TestBuilder_Histogram(t *testing.T) {
	type testHistogram struct {
		name   string
		sum    float64
		count  uint64
		bounds []float64
		bucket []uint64
	}

	goldenHistogram := []testHistogram{
		{
			name:  "cluster.upstream_rq_time",
			sum:   552.5,
			count: 3,
			bounds: []float64{
				0.5,
				1,
				5,
				10,
				25,
				50,
				100,
				250,
				500,
				1000,
				2500,
				5000,
				10000,
				30000,
				60000,
				300000,
				600000,
				1800000,
				3600000,
			},
			bucket: []uint64{
				0,
				0,
				0,
				0,
				0,
				0,
				1,
				2,
				3,
				3,
				3,
				3,
				3,
				3,
				3,
				3,
				3,
				3,
				3,
			},
		},
		{
			name:  "http.downstream_cx_length_ms",
			sum:   562.5,
			count: 3,
			bounds: []float64{
				0.5,
				1,
				5,
				10,
				25,
				50,
				100,
				250,
				500,
				1000,
				2500,
				5000,
				10000,
				30000,
				60000,
				300000,
				600000,
				1800000,
				3600000,
			},
			bucket: []uint64{
				0,
				0,
				0,
				0,
				0,
				0,
				1,
				2,
				3,
				3,
				3,
				3,
				3,
				3,
				3,
				3,
				3,
				3,
				3,
			},
		},
	}
	f := "testdata/histogram"
	labels := map[string]string{
		"name":    uuid.NewString(),
		"cluster": uuid.NewString(),
	}

	bytes, err := os.ReadFile(f)
	must.NoError(t, err)

	histograms := make([]*prompb.MetricFamily, 0)
	must.NoError(t, json.Unmarshal(bytes, &histograms))

	b := NewBuilder(labels)
	for _, histogram := range histograms {
		b.AddHistogram(histogram)
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
	must.Length(t, 2, metricSlice)

	for i := 0; i < metricSlice.Len(); i++ {
		metric := metricSlice.At(0)
		must.Eq(t, pmetric.MetricTypeHistogram, metric.Type(), must.Sprintf("metric types don't match %s!=%s",
			pmetric.MetricTypeHistogram.String(), metric.Type().String()))
	}

	for _, goldenHisto := range goldenHistogram {
		must.Contains[string](t, goldenHisto.name, ContainsMetricName(metricSlice), must.Sprint("Does not contain",
			goldenHisto.name))
		histogram := lookupHistogram(metricSlice, goldenHisto.name)
		must.Eq(t, goldenHisto.sum, histogram.Sum())
		must.Eq(t, goldenHisto.count, histogram.Count())
		must.Eq(t, goldenHisto.bucket, histogram.BucketCounts().AsRaw())
		must.Eq(t, goldenHisto.bounds, histogram.ExplicitBounds().AsRaw())
	}
}

func lookupHistogram(metric pmetric.MetricSlice, name string) pmetric.HistogramDataPoint {
	for i := 0; i < metric.Len(); i++ {
		m := metric.At(i)
		if m.Name() == name {
			return m.Histogram().DataPoints().At(0)
		}
	}
	return pmetric.HistogramDataPoint{}
}

func ContainsMetricName(slice pmetric.MetricSlice) ContainsFunc[string] {
	return func(v string) bool {
		for i := 0; i < slice.Len(); i++ {
			metric := slice.At(i)
			if metric.Name() == v {
				return true
			}
		}
		return false
	}
}
