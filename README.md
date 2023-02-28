<h1>
  <img src="./_doc/logo.svg" align="left" height="46px" alt="Consul logo"/>
  <span>Consul Telemetry Collector</span>
</h1>

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