package prometheus

import (
	"time"

	_go "github.com/prometheus/client_model/go"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

func (b *Builder) AddCounter(family *_go.MetricFamily) {
	// s := sum{}
	otlpMetric := pmetric.NewMetric()

	otlpMetric.SetName(family.GetName())
	otlpMetric.SetDescription(family.GetHelp())
	emptySum := otlpMetric.SetEmptySum()
	emptySum.SetIsMonotonic(true)
	for _, metric := range family.GetMetric() {
		dp := emptySum.DataPoints().AppendEmpty()

		for _, labelPair := range metric.GetLabel() {
			dp.Attributes().PutStr(labelPair.GetName(), labelPair.GetValue())
		}

		t := time.Unix(0, metric.GetTimestampMs()*int64(time.Millisecond))
		dp.SetTimestamp(pcommon.NewTimestampFromTime(t))

		dp.SetDoubleValue(metric.GetCounter().GetValue())
	}

	b.metrics = append(b.metrics, otlpMetric)
}

func NewCounter(family _go.MetricFamily) pmetric.Metric {
	// sum := &otlpmetric.Sum{}
	// https://opentelemetry.io/docs/reference/specification/compatibility/prometheus_and_openmetrics/#counters
	// sum.IsMonotonic = true
	// sum.AggregationTemporality = otlpmetric.AggregationTemporality_AGGREGATION_TEMPORALITY_UNSPECIFIED
	// md := pmetric.NewMetrics()
	// metrics := md.ResourceMetrics().AppendEmpty().ScopeMetrics().AppendEmpty().Metrics().AppendEmpty()
	// metrics.SetName(family.GetName())
	// metrics.SetDescription(family.GetHelp())
	// sum := metrics.Sum().DataPoints()
	//
	// for i, metric := range family.GetMetric() {
	// 	dp := sum.At(i)
	// 	for _, labelPair := range metric.GetLabel() {
	// 		dp.Attributes().PutStr(labelPair.GetName(), labelPair.GetValue())
	// 	}
	// 	dp.SetDoubleValue(metric.GetCounter().GetValue())
	// }

	return pmetric.Metric{}
}
