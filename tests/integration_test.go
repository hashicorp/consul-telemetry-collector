package tests

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	metricsv3 "github.com/envoyproxy/go-control-plane/envoy/service/metrics/v3"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	prom "github.com/prometheus/client_model/go"
	"github.com/shoenig/test/must"
	"github.com/shoenig/test/portal"
	otlpcolmetrics "go.opentelemetry.io/proto/otlp/collector/metrics/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/hashicorp/consul-telemetry-collector/internal/otel"
	"github.com/hashicorp/consul-telemetry-collector/internal/otel/config"
	"github.com/hashicorp/consul-telemetry-collector/internal/otel/config/helpers/exporters"
	"github.com/hashicorp/go-hclog"
)

// impl is an OTLP metrics server.
type impl struct {
	otlpcolmetrics.UnimplementedMetricsServiceServer
	validation func(req *otlpcolmetrics.ExportMetricsServiceRequest)
}

var (
	_ otlpcolmetrics.MetricsServiceServer = &impl{}
)

// RegisterGRPC registers the OTLP metrics server in the grpcServer.
func (s *impl) RegisterGRPC(grpcServer *grpc.Server) {
	otlpcolmetrics.RegisterMetricsServiceServer(grpcServer, s)
}

// RegisterGateway registers this OTLP metrics endpoints in the HTTP web server.
func (s *impl) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	if err := otlpcolmetrics.RegisterMetricsServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		return fmt.Errorf("failed to register metrics handler: %w", err)
	}

	return nil
}

// AdditionalGatewayServeOpts is additional gateway options like matching the resource header.
func (s *impl) AdditionalGatewayServeOpts() []runtime.ServeMuxOption {
	return []runtime.ServeMuxOption{
		runtime.WithMarshalerOption("application/x-protobuf", &runtime.ProtoMarshaller{}),
	}
}

// Export takes in OTLP metrics and writes them to the Prometheus backend.
func (s *impl) Export(ctx context.Context, req *otlpcolmetrics.ExportMetricsServiceRequest) (*otlpcolmetrics.ExportMetricsServiceResponse, error) {
	// hclog.Default().Info("Got resource metrics", "count", len(req.GetResourceMetrics()))
	for _, resourceMetric := range req.GetResourceMetrics() {
		// hclog.Default().Info("Got scope metrics", "count", len(resourceMetric.GetScopeMetrics()))
		for _, scopeMetrics := range resourceMetric.GetScopeMetrics() {
			hclog.Default().Info("Got metrics", "count", len(scopeMetrics.Metrics))
		}
	}
	s.validation(req)
	return &otlpcolmetrics.ExportMetricsServiceResponse{}, nil
}

type Addrs struct {
	GRPCEndpoint string
	HTTPEndpoint string
}

func NewTestServer(t *testing.T, validation func(req *otlpcolmetrics.ExportMetricsServiceRequest)) Addrs {
	t.Helper()
	ctx := context.Background()
	svc := &impl{
		validation: validation,
	}
	svr := grpc.NewServer()
	svc.RegisterGRPC(svr)

	grab := portal.New(t)
	port := grab.One()
	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	must.NoError(t, err)
	go func() {
		must.NoError(t, svr.Serve(ln))
	}()

	t.Cleanup(svr.GracefulStop)

	// GRPC Gateway
	mux := runtime.NewServeMux(svc.AdditionalGatewayServeOpts()...)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	must.NoError(t, svc.RegisterGateway(ctx, mux, ln.Addr().String(), opts))

	f := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hclog.Default().Info("request", "url", r.Header)
		mux.ServeHTTP(w, r)
	})

	httpSvr := httptest.NewServer(f)
	t.Cleanup(httpSvr.Close)
	return Addrs{
		GRPCEndpoint: ln.Addr().String(),
		HTTPEndpoint: httpSvr.Listener.Addr().String(),
	}
}

func Test_OTLPHTTP(t *testing.T) {
	totalMetric := atomic.Int64{}
	addrs := NewTestServer(t, func(req *otlpcolmetrics.ExportMetricsServiceRequest) {
		for _, resourceMetric := range req.GetResourceMetrics() {
			for _, scopeMetric := range resourceMetric.GetScopeMetrics() {
				count := int64(len(scopeMetric.GetMetrics()))
				totalMetric.Add(count)
			}
		}
	})
	envoyPort := portal.New(t).One()
	hclog.Default().Info("Running test server", "addr", addrs)
	collector, err := otel.NewCollector(otel.CollectorCfg{
		ExporterConfig: &config.ExporterConfig{
			ID: exporters.BaseOtlpExporterID,
			Exporter: &exporters.ExporterConfig{
				Endpoint: fmt.Sprintf("http://%s", addrs.HTTPEndpoint),
				Headers: map[string]string{
					"authorization": "abc123",
				},
			},
		},
		MetricsPort:  portal.New(t).One(),
		BatchTimeout: time.Second,
		EnvoyPort:    envoyPort,
	})
	must.NoError(t, err)
	ctx := context.Background()
	go func() { must.NoError(t, collector.Run(ctx)) }()

	total := generateMetrics(t, envoyPort, 30, 30)
	ch := time.After(30 * time.Second)
	for {
		select {
		case <-ch:
			t.Logf("Failed to get %d metrics in 30 seconds", total)
			t.Fail()
		default:
			t.Log(totalMetric.Load(), total)
			time.Sleep(500 * time.Millisecond)
		}
		if totalMetric.Load() == int64(total) {
			break
		}
	}

	collector.Shutdown()
	hclog.Default().Info("Shutting down")
}

func Test_OTLPGRPC(t *testing.T) {
	totalMetric := atomic.Int64{}
	addrs := NewTestServer(t, func(req *otlpcolmetrics.ExportMetricsServiceRequest) {
		for _, resourceMetric := range req.GetResourceMetrics() {
			for _, scopeMetric := range resourceMetric.GetScopeMetrics() {
				count := int64(len(scopeMetric.GetMetrics()))
				totalMetric.Add(count)
			}
		}
	})

	envoyPort := portal.New(t).One()
	hclog.Default().Info("Running test server", "addr", addrs)
	collector, err := otel.NewCollector(otel.CollectorCfg{
		ExporterConfig: &config.ExporterConfig{
			ID: exporters.GRPCOtlpExporterID,
			Exporter: &exporters.ExporterConfig{
				Endpoint: fmt.Sprintf("http://%s", addrs.GRPCEndpoint),
				Headers: map[string]string{
					"authorization": "abc123",
				},
			},
		},
		MetricsPort:  portal.New(t).One(),
		BatchTimeout: time.Second,
		EnvoyPort:    envoyPort,
	})
	must.NoError(t, err)
	ctx := context.Background()
	go func() { must.NoError(t, collector.Run(ctx)) }()

	total := generateMetrics(t, envoyPort, 30, 30)
	ch := time.After(30 * time.Second)
	for {
		select {
		case <-ch:
			t.Logf("Failed to get %d metrics in 30 seconds", total)
			t.Fail()
		default:
			t.Log(totalMetric.Load())
			time.Sleep(500 * time.Millisecond)
		}
		if totalMetric.Load() == int64(total) {
			break
		}
	}

	collector.Shutdown()
	hclog.Default().Info("Shutting down")
}

func ptr[T any](s T) *T {
	return &s
}

const counterNumber int = 10

func generateMetrics(t *testing.T, envoyPort, totalSend, metricCount int) (total int) {
	t.Helper()
	total = totalSend * metricCount
	conn, err := grpc.Dial(fmt.Sprintf("127.0.0.1:%d", envoyPort), grpc.WithBlock(), grpc.WithInsecure())
	must.NoError(t, err)
	client := metricsv3.NewMetricsServiceClient(conn)
	streamClient, err := client.StreamMetrics(context.Background())
	must.NoError(t, err)

	for i := 0; i < totalSend; i++ {
		hclog.Default().Info("sending metric")
		streamClient.Send(&metricsv3.StreamMetricsMessage{
			Identifier: &metricsv3.StreamMetricsMessage_Identifier{
				Node: &corev3.Node{
					Id:      "integration_test",
					Cluster: "test",
				},
			},
			EnvoyMetrics: NewEnvoyMetrics(metricCount),
		})
	}
	return total
}

func NewEnvoyMetrics(metricCount int) []*prom.MetricFamily {
	metrics := make([]*prom.MetricFamily, metricCount)
	for i := 0; i < metricCount; i++ {
		metrics[i] = &prom.MetricFamily{
			Name:   ptr(uuid.NewString()),
			Type:   prom.MetricType_COUNTER.Enum(),
			Metric: NewCounter(counterNumber),
		}
	}
	return metrics
}

func NewCounter(count int) []*prom.Metric {
	metrics := make([]*prom.Metric, count)
	for i := 0; i < count; i++ {
		metrics[i] = &prom.Metric{
			Label: []*prom.LabelPair{},
			Counter: &prom.Counter{
				Value:            ptr(float64(i)),
				CreatedTimestamp: timestamppb.Now(),
			},
			TimestampMs: ptr(time.Now().UnixMilli()),
		}
	}
	return metrics
}
