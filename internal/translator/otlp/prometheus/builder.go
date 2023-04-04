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
		// appendSum(metrics.AppendEmpty(), b.sums[i])
		// s := b.sums[i]
		metric := metrics.AppendEmpty()
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

	// for _, histogram := range b.histograms {

	// }

	metricsDefintion := pmetric.NewMetrics()
	// TODO leave link describing wtf scoped metrics are
	scopedMetrics := metricsDefintion.ResourceMetrics().AppendEmpty().ScopeMetrics().AppendEmpty()
	// append(md.ResourceMetrics().AppendEmpty().ScopeMetrics().AppendEmpty(), metrics)
	b.metricsRef.CopyTo(scopedMetrics.Metrics())

	return metricsDefintion
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
