package confresolver

import (
	"testing"

	"github.com/shoenig/test"
	"github.com/shoenig/test/must"
	"go.opentelemetry.io/collector/component"
)

func Test_NewExtensions(t *testing.T) {
	cfg := &Config{}
	ccfg := cfg.NewExtensions(component.NewID("ballast"))

	ccfg.Set("percentage", 10)

	test.Len(t, 1, cfg.Service.Extensions)
	test.Eq(t, component.NewID("ballast"), cfg.Service.Extensions[0])

	extensions := cfg.Extensions[component.NewID("ballast")]
	extensionCfg, ok := extensions.(componentConfig)
	must.True(t, ok)

	componentConfig, ok := ccfg.(componentConfig)
	must.True(t, ok)

	test.Eq(t, componentConfig["ballast"], extensionCfg["ballast"])
}
