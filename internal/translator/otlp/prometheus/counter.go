package prometheus

import (
	prompb "github.com/prometheus/client_model/go"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

// AddCounter converts a prometheus counter to an OTLP monotonic Sum and adds it to the metrics builder
func (b *Builder) AddCounter(family *prompb.MetricFamily) {
	otlpMetric := pmetric.NewMetric()

	otlpMetric.SetName(normalizeName(family.GetName()))
	otlpMetric.SetDescription(family.GetHelp())
	emptySum := otlpMetric.SetEmptySum()
	emptySum.SetIsMonotonic(true)
	emptySum.SetAggregationTemporality(pmetric.AggregationTemporalityCumulative)
	for _, metric := range family.GetMetric() {
		dp := emptySum.DataPoints().AppendEmpty()

		for _, labelPair := range metric.GetLabel() {
			dp.Attributes().PutStr(labelPair.GetName(), labelPair.GetValue())
		}

		dp.SetTimestamp(timestampFromMs(metric.GetTimestampMs()))
		dp.SetDoubleValue(metric.GetCounter().GetValue())
	}

	b.metrics = append(b.metrics, otlpMetric)
}
