package prometheus

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/google/uuid"
	_go "github.com/prometheus/client_model/go"
	"github.com/shoenig/test/must"
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
	must.Length(t, 1, md.ResourceMetrics().At(0).ScopeMetrics())
	must.Length(t, 2, md.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics())
}
