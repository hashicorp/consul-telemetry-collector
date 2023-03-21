package metrics

import (
	"io"

	metricsv3 "github.com/envoyproxy/go-control-plane/envoy/service/metrics/v3"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/obsreport"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Receiver struct {
	nextConsumer consumer.Metrics
	obsrecv      *obsreport.Receiver
	logger       *zap.Logger
}

var _ metricsv3.MetricsServiceServer = (*Receiver)(nil)

// New creates a new Receiver reference.
func New(nextConsumer consumer.Metrics, logger *zap.Logger) *Receiver {
	logger.Info("Created new receiver")
	return &Receiver{
		nextConsumer: nextConsumer,
		logger:       logger,
	}
}

func (r *Receiver) Register(g *grpc.Server) {
	metricsv3.RegisterMetricsServiceServer(g, r)
}

func (r *Receiver) StreamMetrics(stream metricsv3.MetricsService_StreamMetricsServer) error {

	var identifier *metricsv3.StreamMetricsMessage_Identifier
	for {
		metricsMessage, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return stream.SendAndClose(&metricsv3.StreamMetricsResponse{})
			}
			return err
		}
		if err := metricsMessage.ValidateAll(); err != nil {
			r.logger.Error("failed to validate metric stream", zap.String("error", err.Error()))
			return err
		}

		if identifier != nil {
			identifier = metricsMessage.GetIdentifier()
		}

		metrics := metricsMessage.GetEnvoyMetrics()
		for _, metric := range metrics {
			_ = metric
		}
	}
}
