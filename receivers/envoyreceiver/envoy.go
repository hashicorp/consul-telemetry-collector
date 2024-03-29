// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package envoyreceiver

import (
	"context"
	"errors"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/hashicorp/consul-telemetry-collector/receivers/envoyreceiver/metrics"
)

type envoyReceiver struct {
	cfg             *Config
	logger          *zap.Logger
	grpcServer      *grpc.Server
	metricsReceiver *metrics.Receiver
	shutdownCh      chan struct{}

	settings receiver.CreateSettings
}

var _ receiver.Metrics = (*envoyReceiver)(nil)
var _ component.Component = (*envoyReceiver)(nil)

var _ component.Config = (*Config)(nil)

// Config is the configuration for the envoy receiver.
type Config struct {
	GRPC *configgrpc.GRPCServerSettings `mapstructure:"grpc"`
}

func newEnvoyReceiver(
	set receiver.CreateSettings,
	cfg *Config,
) *envoyReceiver {
	receiver := &envoyReceiver{
		cfg:      cfg,
		settings: set,
		logger: set.TelemetrySettings.Logger.Named("envoyreceiver").With(zap.String("kind", "receiver"),
			zap.String("name", ID)),
	}

	return receiver
}

func (r *envoyReceiver) Start(_ context.Context, host component.Host) error {
	grpcServer, err := r.cfg.GRPC.ToServer(host, r.settings.TelemetrySettings)
	if err != nil {
		r.logger.Error("error creating new server")
		return err
	}

	r.grpcServer = grpcServer
	r.metricsReceiver.Register(grpcServer)

	listener, err := r.cfg.GRPC.ToListener()
	if err != nil {
		return err
	}

	r.logger.Info("Starting GRPC Server", zap.String("endpoint", r.cfg.GRPC.NetAddr.Endpoint))

	r.shutdownCh = make(chan struct{})
	go func() {
		if grpcErr := r.grpcServer.Serve(listener); grpcErr != nil {
			switch {
			case errors.Is(grpcErr, grpc.ErrServerStopped):
				// ignore ErrServerStopped because it's expected
				break
			default:
				host.ReportFatalError(grpcErr)
			}
		}
		r.shutdownCh <- struct{}{}
	}()

	return nil
}

func (r *envoyReceiver) Shutdown(_ context.Context) error {
	if r.grpcServer == nil {
		r.logger.Warn("Shutting down envoy receiver that did not start successfully")
		return nil
	}

	r.logger.Info("Shutting down envoy receiver gracefully")
	r.grpcServer.GracefulStop()
	if r.shutdownCh != nil {
		<-r.shutdownCh
	}

	return nil
}

func (r *envoyReceiver) registerMetrics(nextConsumer consumer.Metrics) {
	r.metricsReceiver = metrics.New(nextConsumer, r.logger)
}
