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

func TestBuilder_Gauge(t *testing.T) {
	type testmetric struct {
		name string
		val  float64
	}

	goldenCounters := []testmetric{
		{
			name: "http.downstream_rq_active",
			val:  0,
		},
		{
			name: "listener_manager.total_listeners_active",
			val:  1,
		},
	}
	f := "testdata/gauge"
	labels := map[string]string{
		"name":    uuid.NewString(),
		"cluster": uuid.NewString(),
	}

	bytes, err := os.ReadFile(f)
	must.NoError(t, err)

	counters := make([]*_go.MetricFamily, 0)
	must.NoError(t, json.Unmarshal(bytes, &counters))

	b := NewBuilder(labels)
	for _, counter := range counters {
		b.AddGauge(counter)
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
		must.Eq(t, pmetric.MetricTypeGauge, metric.Type(), must.Sprintf("metric types don't match %s!=%s",
			pmetric.MetricTypeGauge.String(), metric.Type().String()))
	}

	for _, counter := range goldenCounters {
		must.Contains[string](t, counter.name, ContainsMetricName(metricSlice))
		val := lookupGauge(metricSlice, counter.name)
		must.Eq(t, counter.val, val)
	}
}

func lookupGauge(metric pmetric.MetricSlice, name string) float64 {
	for i := 0; i < metric.Len(); i++ {
		m := metric.At(i)
		if m.Name() == name {
			dp := m.Gauge().DataPoints().At(0).DoubleValue()
			return dp
		}
	}
	return 0
}
