package processors

import (
	"regexp"

	"github.com/hashicorp/go-hclog"
	"go.opentelemetry.io/collector/component"

	"github.com/hashicorp/consul-telemetry-collector/internal/hcp"
)

const filterProcessorName = "filter"

// FilterProcessorID is the component id of the filter processor
var FilterProcessorID component.ID = component.NewID(filterProcessorName)

const (
	regexpMatch = "regexp"
)

// FilterProcessorConfig creates the filter processor configuration
type FilterProcessorConfig struct {
	Metrics *MetricFilters `mapstructure:"metrics"`
}

// MetricFilters is the filter configuration for metrics
type MetricFilters struct {
	Include *MatchProperties `mapstructure:"include"`
}

// MatchProperties specifies how to match against a set of metric names for filtering signals
type MatchProperties struct {
	MatchType   string   `mapstructure:"match_type"`
	MetricNames []string `mapstructure:"metric_names"`
}

// FilterProcessorCfg generates the config for a filter processor
func FilterProcessorCfg(client hcp.TelemetryClient) *FilterProcessorConfig {
	usableFilters := []string{}
	filters, err := client.MetricFilters()
	logger := hclog.Default().Named("config/helpers")
	if err != nil {
		// log failure here, but it's not fatal because the gateway should also filter metrics
		logger.Error("failed to retrieve metric filters from HCP")
	}
	for _, filter := range filters {
		if err := validateFilter(filter); err != nil {
			// log failure here, but it's not fatal because the gateway should also filter metrics
			logger.Error("failed to validate filter", "filter", filter)
			continue
		}
		usableFilters = append(usableFilters, filter)
	}
	logger.Info("created", len(usableFilters), "usable filters")

	cfg := FilterProcessorConfig{
		Metrics: &MetricFilters{
			Include: &MatchProperties{
				MatchType:   regexpMatch,
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
