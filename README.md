<h1>
  <img src="./.doc/logo.svg" align="left" height="46px" alt="Consul logo"/>
  <span>Consul Telemetry Collector</span>
</h1>

Consul Telemetry Collector is a lightweight OpenTelemetry collector used to
collect metrics, logs and traces from a Consul cluster and collect metrics from envoy
service proxies and export them to HCP or another OTLP compliant endpoint. The metric sink
 is encrypted and authorized by the Consul service mesh.

.

Configuration will be loaded in the following order of precedence:

    1. command line opts if specified
    2. env variables if specified
    3. file configuration

# Installation


To install and use the consul-telemetry-collector you will need a Consul version of 1.15.3 or greater and to authorize communication with the collector on the Service Mesh.

## Kubernetes with consul-k8s or Helm

### Table of Contents
* [Minimal Installation](#minimal-installation)
* [Forward metrics to HCP](#forwarding-metrics-to-hcp)

This requires consul-k8s version of 1.1.2 or greater

We currently recommend using the consul-k8s CLI to install the consul-telemetry-collector. If you aren't already using it or the Consul Helm chart you can find instructions and documentation for using it [in the Consul Documentation](https://developer.hashicorp.com/consul/docs/k8s/installation/install). A few small changes to the helm chart enable the telemetry-collector to forward metrics to HCP.

### Minimal Installation

Retrieve the current Helm values from Kubernetes using the `consul-k8s status` command or `helm get values consul` and write them to a file.
```bash
  connectInject:
    enabled: true
  controller:
    enabled: true
  global:
    acls:
      manageSystemACLs: true
    datacenter: mesh-metrics
    name: consul
    tls:
      enableAutoEncrypt: true
      enabled: true
  server:
    affinity: null
    replicas: 3
```

Now add the changes to the helm `values.yaml` file to enable the telemetry-collector deployment.
```yaml
  connectInject:
    enabled: true
  controller:
    enabled: true
  global:
+   metrics:
+     enableTelemetryCollector: true
    acls:
      manageSystemACLs: true
    datacenter: mesh-metrics
    name: consul
    tls:
      enableAutoEncrypt: true
      enabled: true
  server:
    affinity: null
    replicas: 3
+ telemetryCollector:
+    enabled: true
```

Now apply your new configuration using the upgrade command
```bash
consul-k8s upgrade -f values.yaml
```

#### Service Intentions
Next ensure that we authorize communication with the consul-telemetry-collector so that we start receiving envoy metrics.

You'll need to create a `ServiceIntention` to allow that communication.
```bash
cat <<EOF | kubectl apply --filename -
apiVersion: consul.hashicorp.com/v1alpha1
kind: ServiceIntentions
metadata:
  name: global
spec:
  destination:
    name: consul-telemetry-collector
  sources:
  - action: allow
    name: '*'
EOF
```

Now all pods that have been deployed since the Consul upgrade will get an updated envoy configuration and start attempting to forward their metrics to the consul-telemetry-collector.

#### Forward to another otel-collector

To consume these metrics in an APM, forward them to another OTLP metrics HTTP endpoint.

```yaml
...
telemetryCollector:
  enabled: true
  customExporterConfig: |
      {"http_collector_endpoint": "otel-collector:4187"}
```

### Forwarding Metrics to HCP

These metrics can also be sent to HCP's Consul management plane to receive Consul Server and envoy proxy metrics. This assumes that this cluster is already [linked with HCP's Consul management plane](https://developer.hashicorp.com/hcp/docs/consul/usage/management-plane). You'll need the Service Principal and HCP Resource ID for the cluster to authenticate to HCP.

If you've previously used the `cloud` preset to deploy Consul, download the latest version of consul-k8s (>= `1.1.2`) and run `consul-k8s -preset cloud upgrade` to update to the latest version of Consul and enable the telemetry-collector automatically! Skip to [Service Intentions](#service-intentions-1)

If not, follow these instructions to add the new configuration to your values file.

Retrieve the current values.yaml file from Kubernetes using the `consul-k8s status` command or `helm get values consul`
```bash
  connectInject:
    enabled: true
  controller:
    enabled: true
  global:
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
```

Now add the following two snippets to the helm `values.yaml` file to enable the telemetry-collector deployment and upgrade the helm chart.
```yaml
global:
  metrics:
    enableTelemetryCollector: true
```

```yaml
telemetryCollector:
  cloud:
    clientId:
      secretKey: client-id
      secretName: consul-hcp-client-id
    clientSecret:
      secretKey: client-secret
      secretName: consul-hcp-client-secret
  enabled: true
```

After adding the snippets, upgrade the helm chart and you should see a diff that looks like this:

```bash
consul-k8s upgrade -f values.yaml

==> Checking if Consul can be upgraded
 ✓ Existing Consul installation found to be upgraded.
    Name: consul
    Namespace: consul

==> Checking if Consul demo application can be upgraded
    No existing Consul demo application installation found.

==> Consul Upgrade Summary
 ✓ Downloaded charts.

    Difference between user overrides for current and upgraded charts
    -----------------------------------------------------------------

  connectInject:
    enabled: true
  controller:
    enabled: true
  global:
+   metrics:
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
+  cloud:
+    clientId:
+      secretKey: client-id
+      secretName: consul-hcp-client-id
+    clientSecret:
+      secretKey: client-secret
+      secretName: consul-hcp-client-secret
+  enabled: true

    Proceed with upgrade? (Y/n)

```

#### Service Intentions

Next ensure that we authorize communication with the consul-telemetry-collector so that we start receiving envoy metrics.

You'll need to create a ServiceIntention to allow that communication.
```bash
cat <<EOF | kubectl apply --filename -
apiVersion: consul.hashicorp.com/v1alpha1
kind: ServiceIntentions
metadata:
  name: global
spec:
  destination:
    name: consul-telemetry-collector
  sources:
  - action: allow
    name: '*'
EOF
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
