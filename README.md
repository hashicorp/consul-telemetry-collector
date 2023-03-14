<h1>
  <img src="./_doc/logo.svg" align="left" height="46px" alt="Consul logo"/>
  <span>Consul Telemetry Collector</span>
</h1>

Consul Telemetry Collector is a lightweight OpenTelemetry collector used to collect metrics, logs and traces from a Consul cluster and export them to HCP or another OTLP compliant endpoint.


Configuration will be loaded in the following order of precedence:

    1. command line opts if specified
    2. env variables if specified
    3. file configuration


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