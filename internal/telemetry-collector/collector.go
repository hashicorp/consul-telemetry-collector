package telemetrycollector

import (
	"context"
	"log"

	"github.com/hashicorp/consul-telemetry-collector/internal/otelcol"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/loggingexporter"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
	"go.opentelemetry.io/collector/service"
)

func NewCollector() {

}

func Receivers() receiver.Factory {
	return otlpreceiver.NewFactory()
}

func Exporters() exporter.Factory {
	return loggingexporter.NewFactory()
}

func Run(ctx context.Context) error {
	svc, err := otelcol.NewService(ctx, service.Settings{}, service.Config{})
	if err != nil {
		return err
	}

	svc.Logger().Info("hello from svc")

	if err = svc.Start(ctx); err != nil {
		return err
	}
	log.Println("Started otelcol service")

	doneCh := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			err := svc.Shutdown(ctx)
			if err != nil {
				log.Println(err)
			}
			close(doneCh)
		}
	}()
	<-doneCh
	return nil
}
