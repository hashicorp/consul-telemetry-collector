package otelcol

import (
	"context"
	"testing"

	"github.com/shoenig/test"
)

func Test_New(t *testing.T) {
	ctx := context.Background()
	svc, err := New(ctx)
	test.NoError(t, err)
	test.NoError(t, svc.Start(ctx))
	test.NoError(t, svc.Shutdown(ctx))
}
