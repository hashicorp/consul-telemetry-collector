<h1>
  <img src="./.doc/logo.svg" align="left" height="46px" alt="Consul logo"/>
  <span>Consul Telemetry Collector</span>
</h1>

Consul Telemetry Collector is a lightweight OpenTelemetry collector used to
collect metrics, logs and traces from a Consul cluster and collect metrics from envoy
service proxies and export them to HCP or another OTLP compliant endpoint.
The metric sink is encrypted and authorized by the Consul service mesh.

Configuration will be loaded in the following order of precedence:

    1. command line opts if specified
    2. env variables if specified
    3. file configuration

# Installation

To install and use the Consul Telemetry Collector you will need a Consul version of 1.15.3 or greater and to authorize communication with the collector on the Service Mesh.

## Deployment instructions

### Docs

- [Deploy the Consul telemetry collector](https://developer.hashicorp.com/hcp/docs/consul/monitor/management-plane/observability/telemetry-collector)

## Development

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
