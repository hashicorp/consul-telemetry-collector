package metrics

import (
	"io"

	metricsv3 "github.com/envoyproxy/go-control-plane/envoy/service/metrics/v3"
	_go "github.com/prometheus/client_model/go"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/hashicorp/consul-telemetry-collector/internal/translator/otlp/prometheus"
)

// Receiver is the metrics implementation for an envoy metrics receiver
type Receiver struct {
	nextConsumer consumer.Metrics
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

// Register will register the MetricsServiceServer on the provided grpc Server
func (r *Receiver) Register(g *grpc.Server) {
	metricsv3.RegisterMetricsServiceServer(g, r)
}

// StreamMetrics implements the envoy MetricsServiceServer method StreamMetrics.
// It will consume the envoy prometheus metrics and write them to the nextConsumer.
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

		if identifier == nil {
			identifier = metricsMessage.GetIdentifier()
		}

		metrics := metricsMessage.GetEnvoyMetrics()
		// TODO: what are our resource labels from the identifier
		labels := map[string]string{}
		otlpMetrics := translateMetrics(labels, metrics)
		err = r.nextConsumer.ConsumeMetrics(stream.Context(), otlpMetrics)
		if err != nil {
			return err
		}
	}
}

func translateMetrics(resourceLabels map[string]string, envoyMetrics []*_go.MetricFamily) pmetric.Metrics {
	//TODO: use the label provided in
	b := prometheus.NewBuilder(resourceLabels)
	for _, metric := range envoyMetrics {
		switch metric.GetType() {
		case _go.MetricType_COUNTER:
			b.AddCounter(metric)
		}
	}

	return b.Build()
}
