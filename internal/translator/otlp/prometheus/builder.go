package prometheus

import (
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

// Builder is an OTLP metric builder
type Builder struct {
	identity pcommon.Resource
	metrics  []pmetric.Metric
}

// NewBuilder creates a new OTLP metric builder to convert prometheus metrics to OTLP metrics
func NewBuilder(identityLabels map[string]string) *Builder {
	resource := pcommon.NewResource()

	for k, v := range identityLabels {
		resource.Attributes().PutStr(k, v)
	}
	b := &Builder{
		identity: resource,
	}

	return b
}

// Build adds converted metrics to a new pmetric.Metrics
func (b *Builder) Build() pmetric.Metrics {
	metricsSlice := pmetric.NewMetricSlice()

	for _, metric := range b.metrics {
		m := metricsSlice.AppendEmpty()
		metric.CopyTo(m)
	}

	return generateMetricsDefinition(b.identity, metricsSlice)
}

func generateMetricsDefinition(resourceLabels pcommon.Resource, metricsRef pmetric.MetricSlice) pmetric.Metrics {
	// Metrics -> Resource Metrics -> Scope Metrics -> Metrics

	// create the new top level metrics container
	metricsDefintion := pmetric.NewMetrics()

	// create a new resource metrics container
	resourceMetrics := metricsDefintion.ResourceMetrics().AppendEmpty()

	// copy our identity labels onto the resource metrics
	resourceLabels.CopyTo(resourceMetrics.Resource())

	// TODO leave link describing wtf scoped metrics are
	scopedMetrics := resourceMetrics.ScopeMetrics().AppendEmpty()

	// copy our metrics reference into the scope metrics
	metricsRef.CopyTo(scopedMetrics.Metrics())
	return metricsDefintion
}
