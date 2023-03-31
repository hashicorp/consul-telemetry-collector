// Package agent manages the consul-telemetry-collector process, loads the configuration,
// and sets up and manages the lifecycle of the opentelemetry-otel.
package agent

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/consul-telemetry-collector/internal/hcp"
	"github.com/hashicorp/consul-telemetry-collector/internal/otel"
)

// runSvc will initialize and Start the consul-telemetry-collector Service
func runSvc(ctx context.Context, cfg *Config) error {
	logger := hclog.FromContext(ctx)

	var cloud = &Cloud{}
	var hcpClient *hcp.Client
	if cfg.Cloud.IsEnabled() {
		// Set up the HCP SDK here
		logger.Info("Setting up HCP Cloud SDK")
		var err error
		hcpClient, err = hcp.New(cfg.Cloud.ClientID, cfg.Cloud.ClientSecret, cfg.Cloud.ResourceID)
		if err != nil {
			return fmt.Errorf("failed to create hcp client %w", err)
		}
		cloud = cfg.Cloud
	}

	if cfg.HTTPCollectorEndpoint != "" {
		logger.Info("Forwarding telemetry to collector", "addr", cfg.HTTPCollectorEndpoint)
	}

	c := otel.CollectorCfg{
		ClientID:          cloud.ClientID,
		ClientSecret:      cloud.ClientSecret,
		Client:            hcpClient,
		ResourceID:        cloud.ResourceID,
		ForwarderEndpoint: cfg.HTTPCollectorEndpoint,
	}

	collector, err := otel.NewCollector(ctx, c)
	if err != nil {
		return err
	}

	childCtx, cancel := context.WithCancel(ctx)
	go runCollector(childCtx, collector, cancel)
	<-childCtx.Done()
	return nil

}

func runCollector(ctx context.Context, collector otel.Collector, cancel func()) {
	logger := hclog.FromContext(ctx)
	err := collector.Run(ctx)
	if err != nil {
		logger.Error("Failed to run opentelemetry-collector", "error", err)
	}
	cancel()
}
