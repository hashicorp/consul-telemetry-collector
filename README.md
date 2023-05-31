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

We recommend using the Consul helm chart to install the consul-telemetry-collector. If you aren't already using the Consul Helm chart you can find instructions and documentation for using it [in the Consul Documentation](https://developer.hashicorp.com/consul/docs/k8s/installation/install). A few small changes to the helm chart enable the telemetry-collector to forward metrics to HCP. This assumes that this cluster is already [linked with HCP's Consul management plane](https://developer.hashicorp.com/hcp/docs/consul/usage/management-plane).

```yaml
  connectInject:
    enabled: true
  controller:
    enabled: true
  global:
    metrics:
+     enableTelemetryCollector: true
    acls:
      bootstrapToken:
        secretKey: token
        secretName: consul-bootstrap-token
      manageSystemACLs: true
    cloud:
      clientId:
        secretKey: client-id
        secretName: consul-hcp-client-id
      clientSecret:
        secretKey: client-secret
        secretName: consul-hcp-client-secret
      enabled: true
      resourceId:
        secretKey: resource-id
        secretName: consul-hcp-resource-id
    datacenter: mesh-metrics
    gossipEncryption:
      secretKey: key
      secretName: consul-gossip-key
    name: consul
    tls:
      caCert:
        secretKey: tls.crt
        secretName: consul-server-ca
      enableAutoEncrypt: true
      enabled: true
  server:
    affinity: null
    replicas: 3
    serverCert:
      secretName: consul-server-cert
+ telemetryCollector:
+   clientId:
+     secretKey: client-id
+     secretName: consul-hcp-client-id
+   clientSecret:
+     secretKey: client-secret
+     secretName: consul-hcp-client-secret
+   enabled: true
```

Use the custom config to forward metrics to another telemetry-collector.

```yaml
  connectInject:
    enabled: true
  controller:
    enabled: true
  global:
    metrics:
+     enableTelemetryCollector: true
    acls:
      bootstrapToken:
        secretKey: token
        secretName: consul-bootstrap-token
      manageSystemACLs: true
    cloud:
      clientId:
        secretKey: client-id
        secretName: consul-hcp-client-id
      clientSecret:
        secretKey: client-secret
        secretName: consul-hcp-client-secret
      enabled: true
      resourceId:
        secretKey: resource-id
        secretName: consul-hcp-resource-id
    datacenter: mesh-metrics
    gossipEncryption:
      secretKey: key
      secretName: consul-gossip-key
    name: consul
    tls:
      caCert:
        secretKey: tls.crt
        secretName: consul-server-ca
      enableAutoEncrypt: true
      enabled: true
  server:
    affinity: null
    replicas: 3
    serverCert:
      secretName: consul-server-cert
+ telemetryCollector:
+   clientId:
+     secretKey: client-id
+     secretName: consul-hcp-client-id
+   clientSecret:
+     secretKey: client-secret
+     secretName: consul-hcp-client-secret
+   enabled: true
+   customExporterConfig: |
+     {"http_collector_endpoint": "otel-collector:4187"}
```

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
