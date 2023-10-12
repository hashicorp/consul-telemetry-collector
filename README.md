<h1>
  <img src="./.doc/logo.svg" align="left" height="46px" alt="Consul logo"/>
  <span>Consul Telemetry Collector</span>
</h1>

Consul Telemetry Collector is a lightweight [OpenTelemetry Collector](https://opentelemetry.io/docs/collector/) used to collect metrics from Envoy proxies and export them to HCP and other OTLP compliant metrics endpoints. The [Envoy stats sink](https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/metrics/v3/metrics_service.proto#extension-envoy-stat-sinks-metrics-service) is encrypted and authorized by the Consul service mesh.

# Installation

Documentation on deploying the Consul Telemetry Collector to a Kubernetes runtime are provided in the HCP documentation:

- [Deploy the Consul telemetry collector](https://developer.hashicorp.com/hcp/docs/consul/monitor/consul-central/observability/telemetry-collector)

## Configuration

HCP automatically configures these values for you when you link self-managed clusters using `consul-k8s` with the `-cloud` preset. For more information, refer to [Consul telemetry collector in the documentation](https://developer.hashicorp.com/hcp/docs/consul/monitor/consul-central/observability/telemetry-collector).

All configuration options available are listed below. Use the CLIENT_ID, CLIENT_SECRET, and RESOURCE_ID created for observability from HCP Consul Central to send metrics to HCP.

Configurations load in the following order of precedence:

    1. command line opts if specified
    2. env variables if specified
    3. file configuration

```bash
Usage: consul-telemetry-collector agent [options]

        Starts the telemetry-collector and runs until an interrupt is received. The
        collector can forward all metrics to an otlphttp endpoint or to the Hashicorp
        cloud platform.

  -config-file-path=<string>
     Load configuration from a config file.

  -hcp-client-id=<string>
     HCP Service Principal Client ID Environment variable HCP_CLIENT_ID

  -hcp-client-secret=<string>
     HCP Service Principal Client Secret Environment variable
     HCP_CLIENT_SECRET

  -hcp-resource-id=<string>
     HCP Resource ID Environment variable HCP_RESOURCE_ID

  -http-collector-endpoint=<string>
     OTLP HTTP endpoint to forward telemetry to Environment variable
     CO_OTEL_HTTP_ENDPOINT
```

# Development

## Build

### Binary

```bash
make dev
```

### Docker Image

```bash
make docker
```

## Test

### Unit Tests

```bash
make unit-tests
```
