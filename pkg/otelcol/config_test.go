package otelcol

import (
	"context"
	"testing"

	"github.com/shoenig/test"
)

func Test_newConfigProvider(t *testing.T) {
	provider, err := newConfigProvider()
	test.NoError(t, err)

	ctx := context.Background()
	factories, err := components()
	test.NoError(t, err)
	cfg, err := provider.Get(ctx, factories)

	test.NoError(t, err)
	_ = cfg
}
