package processors

import (
	"regexp"

	"go.opentelemetry.io/collector/component"

	"github.com/hashicorp/consul-telemetry-collector/internal/hcp"
)

const filterProcessorName = "filter"

// BatchProcessorID is the component id of the batch processor
var FitlerProcessorID component.ID = component.NewID(filterProcessorName)

const (
	Regexp = "regexp"
)

type FilterProcessorConfig struct {
	Metrics *MetricFilters `mapstructure:"metrics"`
}

type MetricFilters struct {
	Include *MatchProperties `mapstructure:"include"`
}

// MatchProperties specifies the set of properties in a metric to match against and the
// type of string pattern matching to use.
type MatchProperties struct {
	MatchType   string   `mapstructure:"match_type"`
	MetricNames []string `mapstructure:"metric_names"`
}

// FilterProcesorCfg generates the config for a filter processor
func FilterProcessorCfg(client hcp.TelemetryClient) *FilterProcessorConfig {
	usableFilters := []string{}
	filters, err := client.MetricFilters()
	if err != nil {
		// log failure here, but it's not fatal because the gateway should also filter metrics
	}
	for _, filter := range filters {
		if err := validateFilter(filter); err != nil {
			// log failure here, but it's not fatal because the gateway should also filter metrics
			continue
		}
		usableFilters = append(usableFilters, filter)
	}

	cfg := FilterProcessorConfig{
		Metrics: &MetricFilters{
			Include: &MatchProperties{
				MatchType:   Regexp,
				MetricNames: usableFilters,
			},
		},
	}

	return &cfg
}

func validateFilter(s string) error {
	if _, err := regexp.Compile(s); err != nil {
		return err
	}
	return nil
}
