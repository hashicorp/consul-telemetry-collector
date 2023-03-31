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

// Service manages the consul-telemetry-otel. It should be initialized and started by Run
type Service struct {
	// ctx is our lifecycle handler for the Service.
	// All other lifecycle context cancelers should come from this context
	collector otel.Collector
	hcpClient *hcp.Client
}

// runSvc will initialize and Start the consul-telemetry-collector Service
func runSvc(ctx context.Context, cfg *Config) error {
	logger := hclog.Default()

	svc := &Service{}

	var cloud = &Cloud{}
	if cfg.Cloud.IsEnabled() {
		// Set up the HCP SDK here
		logger.Info("Setting up HCP Cloud SDK")
		hcpClient, err := hcp.New(cfg.Cloud.ClientID, cfg.Cloud.ClientSecret, cfg.Cloud.ResourceID)
		if err != nil {
			return fmt.Errorf("failed to create hcp client %w", err)
		}
		svc.hcpClient = hcpClient
		cloud = cfg.Cloud
	}

	if cfg.HTTPCollectorEndpoint != "" {
		logger.Info("Forwarding telemetry to collector", "addr", cfg.HTTPCollectorEndpoint)
	}

	ctx = hclog.WithContext(ctx, logger)

	collector, err := otel.NewCollector(ctx, cfg.HTTPCollectorEndpoint,
		otel.WithCloud(cloud.ResourceID, cloud.ClientID, cloud.ClientSecret, svc.hcpClient),
	)
	if err != nil {
		return err
	}

	svc.collector = collector

	return svc.start(ctx)
}

// Start starts the initialized Service
func (s *Service) start(ctx context.Context) error {
	logger := hclog.FromContext(ctx)
	logger.Info("Shutting down service")
	// We would start the otel collector here
	childCtx, cancel := context.WithCancel(ctx)
	go func() {
		err := s.collector.Run(childCtx)
		hclog.FromContext(ctx).Error("Failed to run opentelemetry-collector", "error", err)
		cancel()
	}()
	<-childCtx.Done()
	s.stop()
	return nil
}

// stop stops a started Service
func (s *Service) stop() {
	s.collector.Shutdown()
}