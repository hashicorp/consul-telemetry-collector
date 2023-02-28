package collector

import (
	"context"

	"github.com/hashicorp/go-hclog"
)

type Service struct {
}

func Run(ctx context.Context, cfg Config) error {
	logger := hclog.Default()
	if cfg.Cloud.IsEnabled() {
		// Set up the HCP SDK here
		logger.Info("Setting up HCP Cloud SDK")
	}

	if cfg.HTTPCollectorEndpoint != "" {
		logger.Info("Forwarding telemetry to collector", "addr", cfg.HTTPCollectorEndpoint)
	}

	svc := new(Service)
	go func() {
		<-ctx.Done()
		logger.Info("Shutting down service")
		svc.Stop()
	}()

	return svc.Start(ctx)
}

func (s *Service) Start(ctx context.Context) error {
	// We would start the otel collector here
	return nil
}

func (s *Service) Stop() {

}
