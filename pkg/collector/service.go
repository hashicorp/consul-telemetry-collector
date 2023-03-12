// Package collector manages the consul-telemetry-collector process, loads the configuration,
// and sets up and manages the lifecycle of the opentelemetry-collector.
package collector

import (
	"context"

	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/consul-telemetry-collector/pkg/otelcol"
)

// Service manages the consul-telemetry-collector. It should be initialized and started by Run
type Service struct {
	// ctx is our lifecycle handler for the Service.
	// All other lifecycle context cancelers should come from this context
	ctx       context.Context
	collector otelcol.Collector
}

// runSvc will initialize and Start the consul-telemetry-collector Service
func runSvc(ctx context.Context, cfg *Config) error {
	logger := hclog.Default()
	if cfg.Cloud.IsEnabled() {
		// Set up the HCP SDK here
		logger.Info("Setting up HCP Cloud SDK")
	}

	if *cfg.HTTPCollectorEndpoint != "" {
		logger.Info("Forwarding telemetry to collector", "addr", cfg.HTTPCollectorEndpoint)
	}

	ctx = hclog.WithContext(ctx, logger)

	collector, err := otelcol.New(ctx)
	if err != nil {
		return err
	}
	svc := &Service{
		ctx:       ctx,
		collector: collector,
	}

	return svc.start(ctx)
}

// Start starts the initialized Service
func (s *Service) start(ctx context.Context) error {
	// We would start the otel collector here
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		err := s.collector.Run(ctx)
		hclog.FromContext(ctx).Error("Failed to run opentelemetry-collector", "error", err)
		cancel()
	}()
	<-ctx.Done()
	logger := hclog.FromContext(s.ctx)
	logger.Info("Shutting down service")
	s.stop()
	return nil
}

// stop stops a started Service
func (s *Service) stop() {
	s.collector.Shutdown()
}
