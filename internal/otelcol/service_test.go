package otelcol

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/service"
)

func Test_NewService(t *testing.T) {
	svc, err := NewService(context.Background(), service.Settings{}, service.Config{})
	require.NoError(t, err)
	l := svc.Logger()
	l.Error("err")
}
