package metrics

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	metricsv3 "github.com/envoyproxy/go-control-plane/envoy/service/metrics/v3"
	"github.com/google/uuid"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/shoenig/test/must"
	"github.com/shoenig/test/portal"
	"github.com/xhhuango/json"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestReceiver_StreamMetrics(t *testing.T) {
	metricSink := new(consumertest.MetricsSink)
	receiver := New(metricSink, zap.NewNop())
	port := portal.New(t).One()

	addr := fmt.Sprintf("127.0.0.1:%d", port)
	l, err := net.Listen("tcp", addr)
	must.NoError(t, err)

	s := grpc.NewServer()
	receiver.Register(s)

	errCh := make(chan error)
	go func() {
		err := s.Serve(l)
		// the grpc.ErrServerStopped is acceptable so send nil over the channel.
		if !errors.Is(err, grpc.ErrServerStopped) {
			errCh <- nil
			return
		}
		errCh <- err
	}()

	// WithBlock() should make sure that we have a connection before calling StreamMetrics().
	// We have an open TCP connection even if the server might not be serving yet, so we shouldn't have a fatal error
	// here.
	conn, err := grpc.Dial(addr,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	must.NoError(t, err)
	client := metricsv3.NewMetricsServiceClient(conn)
	stream, err := client.StreamMetrics(context.Background())
	must.NoError(t, err)

	b, err := os.ReadFile("testdata/source.json")
	must.NoError(t, err)
	envoyMetrics := make([]*io_prometheus_client.MetricFamily, 0)
	must.NoError(t, json.Unmarshal(b, &envoyMetrics))

	err = stream.Send(&metricsv3.StreamMetricsMessage{
		Identifier: &metricsv3.StreamMetricsMessage_Identifier{
			Node: &corev3.Node{
				Id:                   uuid.NewString(),
				Cluster:              "",
				Metadata:             nil,
				DynamicParameters:    nil,
				Locality:             nil,
				UserAgentName:        "",
				UserAgentVersionType: nil,
				Extensions:           nil,
				ClientFeatures:       nil,
			},
		},
		EnvoyMetrics: envoyMetrics,
	})

	must.NoError(t, err)
	must.NoError(t, stream.CloseSend())
	s.GracefulStop()
	// This will check the error from the grpc.Serve
	must.NoError(t, <-errCh)

	// We should have 1 resource metric

	goldenBytes, err := os.ReadFile("testdata/golden.json")
	must.NoError(t, err)
	goldenMetrics, err := new(pmetric.JSONUnmarshaler).UnmarshalMetrics(goldenBytes)
	must.NoError(t, err)

	allmetrics := metricSink.AllMetrics()
	must.Len(t, 1, allmetrics)
	testMetrics := verifyPMetrics(t, allmetrics[0])
	goldenMetricSlice := verifyPMetrics(t, goldenMetrics)

	for i := 0; i < testMetrics.Len(); i++ {
		must.Contains[pmetric.Metric](t, testMetrics.At(i), ContainsMetric(goldenMetricSlice),
			must.Sprintf("metric %d is missing %s", i, spew.Sdump(testMetrics.At(i))))
	}
}

type ContainsFunc[T any] func(T) bool

func (c ContainsFunc[T]) Contains(v T) bool {
	return (c)(v)
}

// ContainsMetric returns a ContainsFunc that looks through a pmetric.MetricSlice to see if the name,
// attributes and value match. It expects each datapoint slice to have 1 value
func ContainsMetric(ms pmetric.MetricSlice) ContainsFunc[pmetric.Metric] {
	return func(needle pmetric.Metric) bool {
		for i := 0; i < ms.Len(); i++ {
			m := ms.At(i)
			if m.Name() != needle.Name() {
				continue
			}
			if m.Type() != needle.Type() {
				continue
			}
			switch needle.Type() {
			case pmetric.MetricTypeSum:
				if m.Sum().AggregationTemporality() != needle.Sum().AggregationTemporality() {
					continue
				}
				if m.Sum().DataPoints().Len() != needle.Sum().DataPoints().Len() {
					continue
				}
				if m.Sum().DataPoints().At(0).DoubleValue() != needle.Sum().DataPoints().At(0).DoubleValue() {
					continue
				}
				if m.Sum().DataPoints().At(0).Attributes().Len() != needle.Sum().DataPoints().At(0).Attributes().Len() {
					continue
				}
			case pmetric.MetricTypeHistogram:
				if m.Histogram().AggregationTemporality() != needle.Histogram().AggregationTemporality() {
					continue
				}
				if m.Histogram().DataPoints().Len() != needle.Histogram().DataPoints().Len() {
					continue
				}
				if m.Histogram().DataPoints().At(0).ExplicitBounds().Len() != needle.Histogram().DataPoints().At(0).ExplicitBounds().Len() {
					continue
				}
				if m.Histogram().DataPoints().At(0).BucketCounts() != needle.Histogram().DataPoints().At(0).BucketCounts() {
					continue
				}
				if m.Histogram().DataPoints().At(0).Attributes().Len() != needle.Histogram().DataPoints().At(0).Attributes().Len() {
					continue
				}
			case pmetric.MetricTypeGauge:
				if m.Gauge().DataPoints().Len() != needle.Gauge().DataPoints().Len() {
					continue
				}
				if m.Gauge().DataPoints().At(0).DoubleValue() != needle.Gauge().DataPoints().At(0).DoubleValue() {
					continue
				}
			}
			return true
		}
		return false
	}
}

func verifyPMetrics(t *testing.T, metrics pmetric.Metrics) pmetric.MetricSlice {
	must.Eq(t, 1, metrics.ResourceMetrics().Len())
	must.Eq(t, 1, metrics.ResourceMetrics().At(0).ScopeMetrics().Len())

	return metrics.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics()
}
