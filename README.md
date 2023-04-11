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

## Usage

To get started run `consul-telemetry-collector agent -dev`. The collector
will run in dev mode and write all collected metrics to the console.

The collector can also:

1. Forward to HCP
2. Forward to another collector

### Forward to HCP

Use the CLIENT_ID, CLIENT_SECRET, and RESOURCE_ID created for observability
from HCP Consul as configuration for. The collector will
then use them to send metrics to HCP.
```
HCP_CLIENT_ID=<client-id> HCP_CLIENT_SECRET=<client-secret>
HCP_RESOURCE_ID=<resource-id> consul-telemetry-collector agent
```

## Development

### Build

#### Binary

```
make dev
```

#### Docker Image

```
make docker
```

### Testing

#### Unit Tests

```
make unit-tests
```
