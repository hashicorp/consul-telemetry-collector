# Consul-k8s or Helm
This requires consul-k8s version of 1.1.2 or greater.

We currently recommend using the consul-k8s CLI to install the Consul Telemetry Collector. If you are not already using the consul-k8s CLI, nor the Consul Helm chart, you can find instructions and documentation [in the Consul Documentation](https://developer.hashicorp.com/consul/docs/k8s/installation/install)

## Step 0
Retrieve the current Helm values from Kubernetes using the `consul-k8s status` command or `helm get values consul` and write them to a file, for example `values.yaml`.
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

## Step 1
Add the changes to the helm `values.yaml` file to enable the Consul Telemetry Collector deployment.
```diff
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

## Step 2
Now, apply your new configuration using the upgrade command:

```bash
consul-k8s upgrade -f values.yaml
```

```bash
helm upgrade consul hashicorp/consul --namespace consul --values values.yaml
```

## Step 3: Service Intentions
Communicate with the Consul Telemetry Collector is protect with Consul's ServiceIntentions. We must authorize communication to start receiving envoy metrics.

Create a `ServiceIntention` to allow that communication:
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

Now, all pods that have been deployed since the Consul upgrade will get an updated envoy configuration and attempt to forward their metrics to the Consul Telemetry Collector.

## Step 4: Forward to another OTEL collector

To consume these metrics in an APM, forward them to another OTLP metrics HTTP endpoint.

```yaml
...
telemetryCollector:
  enabled: true
  customExporterConfig: |
      {"http_collector_endpoint": "otel-collector:4187"}
```
