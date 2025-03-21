module github.com/hashicorp/consul-telemetry-collector/receivers/envoyreceiver

go 1.24

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc
	github.com/envoyproxy/go-control-plane v0.11.1
	github.com/google/uuid v1.3.1
	github.com/hashicorp/consul-telemetry-collector v0.0.0-00010101000000-000000000000
	github.com/prometheus/client_model v0.5.0
	github.com/shoenig/test v0.6.6
	github.com/stretchr/testify v1.8.4
	github.com/xhhuango/json v1.19.0
	go.opentelemetry.io/collector/component v0.88.0
	go.opentelemetry.io/collector/config/configgrpc v0.88.0
	go.opentelemetry.io/collector/config/confignet v0.88.0
	go.opentelemetry.io/collector/confmap v0.88.0
	go.opentelemetry.io/collector/consumer v0.88.0
	go.opentelemetry.io/collector/pdata v1.0.0-rcv0017
	go.opentelemetry.io/collector/receiver v0.88.0
	go.uber.org/zap v1.26.0
	google.golang.org/grpc v1.59.0
)

require (
	github.com/cncf/xds/go v0.0.0-20230607035331-e9ce68804cb4 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.0.2 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.17.1 // indirect
	github.com/knadh/koanf v1.5.0 // indirect
	github.com/knadh/koanf/v2 v2.0.1 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/mapstructure v1.5.1-0.20220423185008-bf980b35cac4 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mostynb/go-grpc-compression v1.2.2 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	go.opentelemetry.io/collector v0.88.0 // indirect
	go.opentelemetry.io/collector/config/configauth v0.88.0 // indirect
	go.opentelemetry.io/collector/config/configcompression v0.88.0 // indirect
	go.opentelemetry.io/collector/config/configopaque v0.88.0 // indirect
	go.opentelemetry.io/collector/config/configtelemetry v0.88.0 // indirect
	go.opentelemetry.io/collector/config/configtls v0.88.0 // indirect
	go.opentelemetry.io/collector/config/internal v0.88.0 // indirect
	go.opentelemetry.io/collector/extension v0.88.0 // indirect
	go.opentelemetry.io/collector/extension/auth v0.88.0 // indirect
	go.opentelemetry.io/collector/featuregate v1.0.0-rcv0017 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.45.0 // indirect
	go.opentelemetry.io/otel v1.19.0 // indirect
	go.opentelemetry.io/otel/metric v1.19.0 // indirect
	go.opentelemetry.io/otel/trace v1.19.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.37.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231012201019-e917dd12ba7a // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/hashicorp/consul-telemetry-collector => ../../
