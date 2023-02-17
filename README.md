# Consul Telemetry Collector

Consul Telemetry Collector is a lightweight openteletry collector used to collect metrics, logs and traces from a Consul cluster and export it to HCP or another OTLP compliant endpoint.

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