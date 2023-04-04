package prometheus

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/google/uuid"
	_go "github.com/prometheus/client_model/go"
	"github.com/shoenig/test/must"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

func TestBuilder_Counter(t *testing.T) {
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
		b.Counter(counter)
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
	must.Length(t, 2, md.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics())
}
