package confhelper

import (
	"go.opentelemetry.io/collector/component"

	"github.com/hashicorp/consul-telemetry-collector/pkg/confresolver"
)

func Ballast(c *confresolver.Config) confresolver.ComponentConfig {
	ballast := c.NewExtensions(component.NewID("memory_ballast"))
	ballast.Set("size_in_percentage", 10)
	return ballast
}
