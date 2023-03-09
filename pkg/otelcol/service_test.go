package otelcol

import (
	"context"
	"testing"
	"time"

	"github.com/shoenig/test"
	"go.opentelemetry.io/collector/otelcol"
)

func Test_New(t *testing.T) {
	ctx := context.Background()
	svc, err := New(ctx)
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
	svc.Shutdown()
}
