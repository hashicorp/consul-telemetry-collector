<h1>
  <img src="./.doc/logo.svg" align="left" height="46px" alt="Consul logo"/>
  <span>Consul Telemetry Collector</span>
</h1>

Consul Telemetry Collector is a lightweight OpenTelemetry collector used to
collect metrics, logs and traces from a Consul cluster and envoy service
mesh and export them to HCP or another OTLP compliant endpoint.

Configuration will be loaded in the following order of precedence:

    1. command line opts if specified
    2. env variables if specified
    3. file configuration

## Installation

### Kubernetes

We recommend using the Consul helm chart to install the consul-telemetry-collector. If you aren't already using the Consul Helm chart you can find instructions and documentation for using it [in the Consul Documentation](https://developer.hashicorp.com/consul/docs/k8s/installation/install) A few small changes to the helm chart are all that are necessary to enable the telemetry-collector and forward metrics to HCP.

```yaml
globals:
  metrics:
    enableTelemetryCollector: true
telemetryCollector:
  enabled: true
  cloud:
    clientId: # These should match the Kubernetes Secret's for the HCP ClientID and HCP ClientSecret
      secretName: hcp-client-id
      secretKey: client-id
    clientSecret:
      secretName: hcp-client-secret
      secretKey: client-secret
```

Use the custom config to forward metrics to another telemetry-collector.

```yaml
globals:
  metrics:
    enableTelemetryCollector: true
telemetryCollector:
  enabled: true
  customExporterConfig: |
    {"http_collector_endpoint": "otel-collector:4187"
```

## Usage

To get started run `consul-telemetry-collector agent -dev`. The collector
will run in dev mode and write all collected metrics to the console.

The collector can also:

1. Forward to HCP
2. Forward to another collector

### Forward to HCP

Use the CLIENT_ID, CLIENT_SECRET, and RESOURCE_ID created for observability
from HCP Consul. The collector will use them to send metrics to HCP.

```bash
HCP_CLIENT_ID=<client_id> HCP_CLIENT_SECRET=<client_secret> HCP_RESOURCE_ID=<resource_id> consul-telemetry-collector agent
```

## Development

### Build

#### Binary

```bash
make dev
```

#### Docker Image

```bash
make docker
```

### Testing

#### Unit Tests

```bash
make unit-tests
```
