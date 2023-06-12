// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package prometheus

import (
	prompb "github.com/prometheus/client_model/go"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

// AddGauge converts a prometheus gauge to an OTLP Gauge and adds it to the metrics builder.
func (b *Builder) AddGauge(family *prompb.MetricFamily) {
	otlpMetric := pmetric.NewMetric()

	otlpMetric.SetName(normalizeName(family.GetName()))
	otlpMetric.SetDescription(family.GetHelp())
	emptyGauge := otlpMetric.SetEmptyGauge()
	for _, metric := range family.GetMetric() {
		dp := emptyGauge.DataPoints().AppendEmpty()

		for _, labelPair := range metric.GetLabel() {
			dp.Attributes().PutStr(labelPair.GetName(), labelPair.GetValue())
		}

		dp.SetTimestamp(timestampFromMs(metric.GetTimestampMs()))
		dp.SetDoubleValue(metric.GetGauge().GetValue())
	}

	b.metrics = append(b.metrics, otlpMetric)
}
