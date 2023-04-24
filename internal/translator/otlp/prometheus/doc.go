// Package prometheus implements a translator to convert prometheus metrics to OTLP metrics.
// The translation is expected to work with the envoy metricsserver which emits all metrics
// as prometheus protobufs. Counters should be cumulative and only gauges, counters and
// histograms are translated.
//
// Histograms that are emitted by the envoy metrics server are delta histograms instead of cumulative
package prometheus
