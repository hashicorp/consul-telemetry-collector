package tests

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hashicorp/go-hclog"
	"github.com/shoenig/test/must"
	"github.com/shoenig/test/portal"
	otlpcolmetrics "go.opentelemetry.io/proto/otlp/collector/metrics/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
