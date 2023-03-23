package otelcol

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shoenig/test"
	"github.com/shoenig/test/must"
	"go.opentelemetry.io/collector/featuregate"
	"go.opentelemetry.io/collector/otelcol"

	"github.com/hashicorp/consul-telemetry-collector/internal/hcp"
	"github.com/hashicorp/hcp-sdk-go/resource"
)

func Test_New(t *testing.T) {
	testcases := map[string]struct {
		rid  string
		cid  string
		csec string
	}{
		"WithoutHCP": {
			rid:  "",
			cid:  "",
			csec: "",
		},
		"WithHCP": {
			rid: resource.Resource{
				ID:           "cluster",
				Type:         "type",
				Organization: uuid.NewString(),
				Project:      uuid.NewString(),
			}.String(),
			cid:  uuid.NewString(),
			csec: uuid.NewString(),
		},
	}
	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			r := tc.rid

			var client *hcp.MockClient
			if r != "" {
				client = &hcp.MockClient{
					MockMetricsEndpoint: "https://metrics.test.endpoint",
				}
			}

			ctx := context.Background()

			svc, err := New(ctx, "https://localhost:4138", tc.rid, tc.cid, tc.csec, client)
			test.NoError(t, err)
			go func() {
				err := svc.Run(ctx)
				test.NoError(t, err)
			}()

			timer := time.NewTimer(1 * time.Second)

			for {
				select {
				case <-timer.C:
					t.Fail()
					break
				default:
				}
				if svc.GetState() == otelcol.StateRunning {
					break
				}
			}

			timer.Stop()

			features := make([]*featuregate.Gate, 0)
			featuregate.GlobalRegistry().VisitAll(func(gate *featuregate.Gate) {
				features = append(features, gate)
			})

			var otelGate *featuregate.Gate
			test.Contains[string](t, otelFeatureGate, containsFunc[string](func(s string) bool {
				for _, gate := range features {
					if gate.ID() == otelFeatureGate {
						otelGate = gate
						return true
					}
				}
				return false
			}))

			must.NotNil(t, otelGate)
			test.True(t, otelGate.IsEnabled())
			svc.Shutdown()
		})
	}

}

type containsFunc[T any] func(T) bool

func (c containsFunc[T]) Contains(s T) bool {
	return c(s)
}
