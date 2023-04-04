package prometheus

import (
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

type Builder struct {
	identity pcommon.Resource
	id       map[string]string
	sums     []sum
	total    int

	// Metrics -> Resource Metrics -> Scope Metrics -> Metrics
	// md is a reference to the
	md pmetric.Metrics
	// metrics is a reference to underlying data points
	metricsRef pmetric.MetricSlice
}

func NewBuilder(identityLabels map[string]string) *Builder {
	resource := pcommon.NewResource()

	for k, v := range identityLabels {
		resource.Attributes().PutStr(k, v)
	}
	b := &Builder{
		identity: resource,
		id:       identityLabels,
	}

	// TODO: Add resource labels to the Resource under ResourceMetrics
	b.metricsRef = pmetric.NewMetricSlice()

	return b
}

func (b *Builder) Build() pmetric.Metrics {
	metrics := b.metricsRef

	for _, s := range b.sums {
		appendSum(metrics.AppendEmpty(), s)
	}

	// for _, histogram := range b.histograms {

	// }

	return generateMetricsDefinition(b.identity, metrics)
}

func appendSum(metric pmetric.Metric, s sum) {
	metric.SetName(s.name)
	metric.SetDescription(s.help)

	dpSlice := metric.SetEmptySum().DataPoints()
	for _, v := range s.value {
		dp := dpSlice.AppendEmpty()
		dp.SetDoubleValue(v.value)
		dp.SetTimestamp(pcommon.NewTimestampFromTime(s.timestamp))
		for k, v := range v.label {
			dp.Attributes().PutStr(k, v)
		}
	}
}

func generateMetricsDefinition(resourceLabels pcommon.Resource, metricsRef pmetric.MetricSlice) pmetric.Metrics {
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
