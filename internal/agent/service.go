// Package agent manages the consul-telemetry-collector process, loads the configuration,
// and sets up and manages the lifecycle of the opentelemetry-otel.
package agent

import (
	"context"
	"fmt"

	"github.com/hashicorp/consul-telemetry-collector/internal/hcp"
	"github.com/hashicorp/consul-telemetry-collector/internal/otel"
	"github.com/hashicorp/go-hclog"
)

// Service runs a otel.Collector with a configured otel pipeline.
type Service struct {
	cfg       otel.CollectorCfg
	collector otel.Collector
}

// NewService returns a new Service based off the past in configuration.
func NewService(cfg *Config) (*Service, error) {
	s := &Service{}
	s.cfg = otel.CollectorCfg{ForwarderEndpoint: cfg.HTTPCollectorEndpoint}

	if cfg.Cloud != nil && cfg.Cloud.IsEnabled() {
		hcpClient, err := hcp.New(&hcp.Params{
			ClientID:     cfg.Cloud.ClientID,
			ClientSecret: cfg.Cloud.ClientSecret,
			ResourceURL:  cfg.Cloud.ResourceID,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create hcp client %w", err)
		}
		s.cfg.ClientID = cfg.Cloud.ClientID
		s.cfg.ClientSecret = cfg.Cloud.ClientSecret
		s.cfg.Client = hcpClient
		s.cfg.ResourceID = cfg.Cloud.ResourceID
	}
	var err error
	s.collector, err = otel.NewCollector(s.cfg)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// Run will initialize and Start the consul-telemetry-collector Service.
func (s *Service) Run(ctx context.Context) error {
	logger := hclog.FromContext(ctx)

	go s.handleShutdown(ctx)

	// blocking call
	err := s.collector.Run(ctx)
	if err != nil {
		logger.Error("failed to run opentelemetry-collector", "error", err)
		return err
	}
	return nil
}

func (s *Service) handleShutdown(ctx context.Context) {
	<-ctx.Done()
	s.collector.Shutdown()
}
