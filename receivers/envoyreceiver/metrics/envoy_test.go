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

	b, err := os.ReadFile("testdata/metric.source.1.golden")
	must.NoError(t, err)
	envoyMetrics := make([]*io_prometheus_client.MetricFamily, 0)
	must.NoError(t, json.Unmarshal(b, &envoyMetrics))

	// spew.Dump(envoyMetrics)

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
	stream.CloseSend()
	s.GracefulStop()
	// This will check the error from the grpc.Serve
	must.NoError(t, <-errCh)

	// We should have 1 resource metric
	must.Len(t, 1, metricSink.AllMetrics())
	must.Eq(t, 1, metricSink.AllMetrics()[0].ResourceMetrics().Len())
	must.Eq(t, 1, metricSink.AllMetrics()[0].ResourceMetrics().At(0).ScopeMetrics().Len())

	scopeMetrics := metricSink.AllMetrics()[0].ResourceMetrics().At(0).ScopeMetrics()

	for i := 0; i < scopeMetrics.At(0).Metrics().Len(); i++ {
		metrics := scopeMetrics.At(0).Metrics().At(i)
		spew.Dump(metrics)
	}

}
