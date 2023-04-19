package prometheus

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/google/uuid"
	_go "github.com/prometheus/client_model/go"
	"github.com/shoenig/test/must"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

func TestBuilder_Counter(t *testing.T) {
	type testcounter struct {
		name string
		val  float64
	}

	goldenCounters := []testcounter{
		{
			name: "cluster.upstream_cx_rx_bytes",
			val:  85447,
		},
		{
			name: "cluster.update_attempt",
			val:  4,
		},
	}
	f := "testdata/counter"
	labels := map[string]string{
		"name":    uuid.NewString(),
		"cluster": uuid.NewString(),
	}

	counterBytes, err := os.ReadFile(f)
	must.NoError(t, err)

	counters := make([]*_go.MetricFamily, 0)
	must.NoError(t, json.Unmarshal(counterBytes, &counters))

	b := NewBuilder(labels)
	for _, counter := range counters {
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
	must.Length(t, 2, metricSlice)

	for i := 0; i < metricSlice.Len(); i++ {
		metric := metricSlice.At(0)
		must.Eq(t, pmetric.MetricTypeSum, metric.Type())
		must.True(t, metric.Sum().IsMonotonic())
	}

	for _, counter := range goldenCounters {
		must.Contains[string](t, counter.name, ContainsMetricName(metricSlice))
		val := lookupSum(metricSlice, counter.name)
		must.Eq(t, counter.val, val)
	}
}

type ContainsFunc[T any] func(T) bool

func (c ContainsFunc[T]) Contains(v T) bool {
	return (c)(v)
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

func lookupSum(metric pmetric.MetricSlice, name string) float64 {
	for i := 0; i < metric.Len(); i++ {
		m := metric.At(i)
		if m.Name() == name {
			dp := m.Sum().DataPoints().At(0).DoubleValue()
			return dp
		}
	}
	return 0
}
