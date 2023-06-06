# Collecting Consul Service Mesh Metrics for HCP Managed Consul Cluster

These instructions describe how to add Envoy metric collection to an existing Consul service mesh in Kubernetes where the Consul cluster is [managed by HCP](https://developer.hashicorp.com/hcp/docs/consul/usage/management-plane).

## Requirements

- Consul version `>= 1.15.3`
- [consul-k8s](https://developer.hashicorp.com/consul/docs/k8s/installation/install-cli#install-the-cli) version `>= 1.1.2` OR the [Consul Helm chart](https://developer.hashicorp.com/consul/docs/k8s/installation/install#install-consul-on-kubernetes-with-helm)
- Consul cluster linked to [HCP's Consul Management Plane](https://developer.hashicorp.com/hcp/docs/consul/usage/management-plane)

## Step 0: Upgrading a `consul-k8s` Deployed Datacenter

If you used the `-preset cloud` of `consul-k8s` to deploy Consul:
1. [download](https://developer.hashicorp.com/consul/docs/k8s/installation/install-cli#install-the-cli) `consul-k8s` version `>= 1.1.2`
1. run `consul-k8s -preset cloud upgrade` to update to the latest version of Consul and enable the Consul Telemetry Collector
1. skip to [Step 2: Configure Service Intentions](#step-2-configure-service-intentions)

## Step 1: Deploy Consul Telemetry Collector

Retrieve the current Helm values from Kubernetes using the `consul-k8s status` command or `helm get values consul` and write them to a file, for example `values.yaml`. The configuration file below is an example and yours may have additional settings:

```bash
connectInject:
  enabled: true
controller:
  enabled: true
global:
  cloud:
    enabled: true
    clientId:
      secretKey: client-id
      secretName: consul-hcp-client-id
    clientSecret:
      secretKey: client-secret
      secretName: consul-hcp-client-secret
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

Add the following settings to the helm `values.yaml` file to enable the Consul Telemetry Collector deployment:

```diff
  connectInject:
    enabled: true
  controller:
    enabled: true
  global:
+   metrics:
+     enableTelemetryCollector: true
    cloud:
      enabled: true
      clientId:
        secretKey: client-id
        secretName: consul-hcp-client-id
      clientSecret:
        secretKey: client-secret
        secretName: consul-hcp-client-secret
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

Apply the new configuration using the upgrade command:

```bash
# if using consul-k8s
consul-k8s upgrade -f values.yaml

# if using helm
helm install consul hashicorp/consul --create-namespace --namespace consul --values values.yaml
```

## Step 2: Configure Service Intentions

You will need to authorize communication to the Consul Telemetry Collector from Envoy proxies with a [`ServiceIntention`](https://developer.hashicorp.com/consul/docs/connect/config-entries/service-intentions):

```bash
cat <<EOF | kubectl apply --filename -
apiVersion: consul.hashicorp.com/v1alpha1
kind: ServiceIntentions
metadata:
  name: consul-telemetry-collector
spec:
  destination:
    name: consul-telemetry-collector
  sources:
  - action: allow
    name: '*'
EOF
```

## Step 3: (Optional) Forward to Another OTEL Collector

To consume these metrics in [another OTLP-compatible collector or back end](https://opentelemetry.io/docs/concepts/sdk-configuration/otlp-exporter-configuration/#otel_exporter_otlp_metrics_endpoint), configure the Consul Telemetry Collector with the `telemetryCollector.customExportConfig` setting:

```yaml
...
telemetryCollector:
  enabled: true
  customExporterConfig: |
      {"http_collector_endpoint": "otel-collector:4187"}
```
