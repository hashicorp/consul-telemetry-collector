# Development

To get started run `consul-telemetry-collector agent -dev`. The collector
will run in dev mode and write all collected metrics to the console.

The collector can also:

1. Forward to HCP
2. Forward to another collector

## Forward to HCP

Use the CLIENT_ID, CLIENT_SECRET, and RESOURCE_ID created for observability
from HCP Consul. The collector will use them to send metrics to HCP.

```bash
HCP_CLIENT_ID=<client_id> \
HCP_CLIENT_SECRET=<client_secret> \
HCP_RESOURCE_ID=<resource_id> \
consul-telemetry-collector agent
```

## Build

### Binary

```bash
make dev
```

### Docker Image

```bash
make docker
```

## Testing

### Unit Tests

```bash
make unit-tests
```
