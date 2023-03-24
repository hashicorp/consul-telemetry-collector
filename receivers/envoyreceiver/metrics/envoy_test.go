package metrics

import (
	"context"
	"fmt"
	"net"
	"testing"

	metricsv3 "github.com/envoyproxy/go-control-plane/envoy/service/metrics/v3"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/shoenig/test/must"
	"github.com/shoenig/test/portal"
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

	go func() {
		err = s.Serve(l)
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

	err = stream.Send(&metricsv3.StreamMetricsMessage{
		Identifier:   &metricsv3.StreamMetricsMessage_Identifier{},
		EnvoyMetrics: []*io_prometheus_client.MetricFamily{},
	})

	must.NoError(t, err)
}
