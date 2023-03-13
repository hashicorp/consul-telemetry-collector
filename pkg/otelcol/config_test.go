package otelcol

import (
	"context"
	"testing"

	"github.com/shoenig/test"
)

func Test_newConfigProvider(t *testing.T) {
	provider, err := newConfigProvider("https://localhost:4138")
	test.NoError(t, err)

	ctx := context.Background()

	// This provider.Get call will perform a configuration retrieval and ensure that it can be unmarshal'd in the
	// expected component config. To perform that Unmarshal we need the actual component code to unmarshal the map
	// [string]interface{} into the receiver/exporter/etc Config struct.
	factories, err := components()
	test.NoError(t, err)
	cfg, err := provider.Get(ctx, factories)

	test.NoError(t, err)
	test.NoError(t, cfg.Validate())

}
