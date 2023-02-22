package otelcol

import (
	"context"
	"testing"

	"go.opentelemetry.io/collector/service"

	"go.opentelemetry.io/collector/component"

	"github.com/stretchr/testify/require"
)

var testcfg string = `
receivers:
  otlp:
    protocols:
      http: {}

service:
  pipelines:
    traces:
      receivers: [otlp]
      exporters: [otlp]
`

func TestStaticCfg(t *testing.T) {
	t.Run("TestStaticConfigResolver", func(t *testing.T) {
		ctx := context.Background()
		cp, err := Provider()
		require.NoError(t, err)

		factories, err := components()
		require.NoError(t, err)
		cfg, err := cp.Get(ctx, factories)

		require.NoError(t, err)

		require.Contains(t, cfg.Receivers, component.NewID("otlp"))
		require.Contains(t, cfg.Exporters, component.NewID("logging"))
	})
}

func TestNewService(t *testing.T) {
	ctx := context.Background()
	svc, err := NewService(ctx, service.Settings{}, service.Config{})
	require.NoError(t, err)
	err = svc.Start(ctx)
	require.NoError(t, err)
	err = svc.Shutdown(ctx)
	require.NoError(t, err)
}
