// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package prometheus

import (
	"math"

	prompb "github.com/prometheus/client_model/go"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

// AddHistogram converts a prometheus histogram to an OTLP histogram.
func (b *Builder) AddHistogram(family *prompb.MetricFamily) {
	otlpMetric := pmetric.NewMetric()

	otlpMetric.SetName(normalizeName(family.GetName()))
	otlpMetric.SetDescription(family.GetHelp())

	emptyHistogram := otlpMetric.SetEmptyHistogram()
	emptyHistogram.SetAggregationTemporality(pmetric.AggregationTemporalityCumulative)
	for _, metric := range family.GetMetric() {
		histogram := metric.GetHistogram()

		if !isValidHistogram(histogram) {
			continue
		}

		dp := emptyHistogram.DataPoints().AppendEmpty()
		dp.SetCount(histogram.GetSampleCount())
		dp.SetSum(histogram.GetSampleSum())

		bounds, bucket := getBoundsAndBuckets(histogram)

		dp.BucketCounts().FromRaw(bucket)
		dp.ExplicitBounds().FromRaw(bounds)

		for _, labelPair := range metric.GetLabel() {
			dp.Attributes().PutStr(labelPair.GetName(), labelPair.GetValue())
		}

		dp.SetTimestamp(timestampFromMs(metric.GetTimestampMs()))
	}

	b.metrics = append(b.metrics, otlpMetric)
}

func getBoundsAndBuckets(histogram *prompb.Histogram) (bounds []float64, bucketCount []uint64) {
	bounds = []float64{}
	bucketCount = []uint64{}

	for _, bucket := range histogram.GetBucket() {
		if math.IsNaN(bucket.GetUpperBound()) {
			continue
		}
		bounds = append(bounds, bucket.GetUpperBound())
		bucketCount = append(bucketCount, bucket.GetCumulativeCount())
	}

	return bounds, bucketCount
}

func isValidHistogram(histogram *prompb.Histogram) bool {
	if histogram.SampleCount == nil || histogram.SampleSum == nil {
		return false
	}

	if len(histogram.GetBucket()) == 0 {
		return false
	}
	return true
}
