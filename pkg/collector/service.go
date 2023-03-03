// Package collector manages the consul-telemetry-collector process, loads the configuration,
// and sets up and manages the lifecycle of the opentelemetry-collector.
package collector

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/v2/hclsimple"
)

// Service manages the consul-telemetry-collector. It should be initialized and started by Run
type Service struct {
}

// Run will initialize and Start the consul-telemetry-collector Service
func Run(ctx context.Context, cfg Config) error {
	logger := hclog.Default()

	if err := LoadConfig(cfg.ConfigFile, &cfg); err != nil {
		return err
	}

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

// Start starts the initialized Service
func (s *Service) Start(ctx context.Context) error {
	// We would start the otel collector here
	return nil
}

// Stop stops a started Service
func (s *Service) Stop() {

}

// LoadConfig will read, and marshal a configuration file with hclsimple and merge it with the provided Config
func LoadConfig(configFile string, cfg *Config) error {
	if configFile == "" {
		return nil
	}
	f, err := os.Open(configFile)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", configFile, err)
	}

	// wrap our file in a 1mb limit reader
	mb := int64(1024 * 1024 * 1024)
	r := io.LimitReader(f, mb)
	buffer := bytes.NewBuffer(make([]byte, 0, mb))
	_, err = io.Copy(buffer, r)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", configFile, err)
	}

	return hclsimple.Decode(configFile, buffer.Bytes(), nil, cfg)
}
