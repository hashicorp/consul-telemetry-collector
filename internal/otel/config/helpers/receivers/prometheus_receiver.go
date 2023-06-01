// Package receivers holds the type of receivers that consul telemetery supports
package receivers

import (
	"go.opentelemetry.io/collector/component"
)

const prometheusReceiverName = "prometheus"
const scrapeConfigKey = "scrape_configs"

// PrometheusReceiverID is the component id of the prometheus receiver.
var PrometheusReceiverID component.ID = component.NewID(prometheusReceiverName)

// PrometheusConfig is the highest container for a prometheus config
type PrometheusConfig struct {
	// We'll probably want to make this a little more generic in the future if we scrape more targets
	Config map[string][]ScrapeConfig `mapstructure:"config"`
}

// ScrapeConfig matches a single minimal scrape configs for prometheus
type ScrapeConfig struct {
	JobName        string         `mapstructure:"job_name"`
	ScrapeInterval string         `mapstructure:"scrape_interval"`
	StaticConfigs  []StaticConfig `mapstructure:"static_configs,omitempty"`
}

// StaticConfig a minimal prometheus static scrape config
type StaticConfig struct {
	Targets []string `mapstructure:"targets"`
}

// PrometheusReceiverCfg  generates the prometheus config for scraping the local telemetry-collector metrics
func PrometheusReceiverCfg() *PrometheusConfig {
	// This should create a config that looks like this for scraping our own metrics
	/*
		prometheus:
		  config:
			scrape_configs:
			- job_name: consul-telemetry-collector
			  scrape_interval: 10s
			  static_configs:
			  - targets:
				- localhost:9090
	*/

	return &PrometheusConfig{
		Config: map[string][]ScrapeConfig{
			scrapeConfigKey: []ScrapeConfig{
				{
					JobName:        "consul-telemetry-collector",
					ScrapeInterval: "10s",
					StaticConfigs: []StaticConfig{
						{
							Targets: []string{"localhost:9090"},
						},
					},
				},
			},
		},
	}
}
