package prometheus

import (
	"time"

	_go "github.com/prometheus/client_model/go"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

// otlp type
type sum struct {
	name        string
	help        string
	isMonotonic bool
	value       []counter
	timestamp   time.Time
}

// underlying metric type from prom
type counter struct {
	label map[string]string
	value float64
}

func (b *Builder) Counter(family *_go.MetricFamily) {
	s := sum{}
	s.name = family.GetName()
	s.help = family.GetHelp()
	var firstTimestampMs int64
	for _, metric := range family.GetMetric() {
		if *metric.TimestampMs < firstTimestampMs {
			firstTimestampMs = *metric.TimestampMs
		}

		c := counter{
			label: map[string]string{},
			value: metric.GetCounter().GetValue(),
		}

		for _, labelPair := range metric.GetLabel() {
			c.label[labelPair.GetName()] = labelPair.GetValue()
		}

		s.value = append(s.value, c)
	}

	s.timestamp = time.Unix(0, firstTimestampMs*int64(time.Millisecond))
	b.sums = append(b.sums, s)
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
