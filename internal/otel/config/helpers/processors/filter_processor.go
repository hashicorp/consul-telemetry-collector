// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package processors

import (
	"fmt"
	"regexp"

	"go.opentelemetry.io/collector/component"

	"github.com/hashicorp/consul-telemetry-collector/internal/hcp"
	"github.com/hashicorp/go-hclog"
)

const filterProcessorName = "filter"

// FilterProcessorID is the component id of the filter processor.
var FilterProcessorID component.ID = component.NewID(filterProcessorName)

const (
	regexpMatchType = "regexp"
)

// FilterProcessorConfig creates the filter processor configuration.
type FilterProcessorConfig struct {
	Metrics *MetricFilters `mapstructure:"metrics"`
}

// MetricFilters is the filter configuration for metrics.
type MetricFilters struct {
	Include *MatchProperties `mapstructure:"include"`
}

// MatchProperties specifies how to match against a set of metric names for filtering signals.
type MatchProperties struct {
	MatchType   string   `mapstructure:"match_type"`
	MetricNames []string `mapstructure:"metric_names"`
}

// FilterProcessorCfg generates the config for a filter processor.
func FilterProcessorCfg(client hcp.TelemetryClient) *FilterProcessorConfig {
	usableFilters := []string{}
	filters, err := client.MetricFilters()
	logger := hclog.Default().Named("config/helpers")
	if err != nil {
		// log failure here, but it's not fatal because the gateway should also filter metrics
		logger.Warn("failed to retrieve metric filters from HCP", "error", err)
		return &FilterProcessorConfig{}
	}
	for _, filter := range filters {
		if err := validateFilter(filter); err != nil {
			// log failure here, but it's not fatal because the gateway should also filter metrics
			logger.Warn("failed to validate filter", "filter", filter, "error", err)
			continue
		}
		usableFilters = append(usableFilters, filter)
	}
	logger.Info(fmt.Sprintf("created %d usable filters for the HCP pipeline", len(usableFilters)))

	cfg := FilterProcessorConfig{
		Metrics: &MetricFilters{
			Include: &MatchProperties{
				MatchType:   regexpMatchType,
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
