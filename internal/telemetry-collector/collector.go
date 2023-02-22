package telemetrycollector

import (
	"context"
	"log"

	"github.com/hashicorp/consul-telemetry-collector/internal/otelcol"
	"go.opentelemetry.io/collector/service"
)

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
