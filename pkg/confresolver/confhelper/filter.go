package confhelper

import (
	"regexp"
	"strings"

	"go.opentelemetry.io/collector/component"

	"github.com/hashicorp/consul-telemetry-collector/pkg/confresolver"
)

// Filter confhelper creates a regex based inclusion filter based on the provided list of filters
func Filter(c *confresolver.Config, includeFilter []string, pipeline confresolver.PipelineIDer,
	pipelines ...confresolver.PipelineIDer) error {
	if len(includeFilter) == 0 {
		return nil
	}

	cfg := c.NewProcessor(component.NewID("filter"), pipeline, pipelines...)

	include := cfg.SetMap("metrics").SetMap("include")
	include.Set("match_type", "regexp")

	for _, filter := range includeFilter {
		_, err := regexp.Compile(filter)
		if err != nil {
			return err
		}
	}

	include.Set("metric_names", strings.Join(includeFilter, ","))
	return nil
}
