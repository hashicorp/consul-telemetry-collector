package confresolver

import (
	"testing"

	"github.com/shoenig/test"
	"go.opentelemetry.io/collector/component"
)

func Test_Pipeline(t *testing.T) {
	c := &Config{}
	metricsPipelineRef := c.NewPipeline(component.DataTypeMetrics)
	tracesPipelineRef := c.NewPipeline(component.DataTypeTraces)
	pipelines := c.Service.Pipelines

	metricPipeline, ok := pipelines[metricsPipelineRef.id()]
	test.NotNil(t, metricPipeline)
	test.True(t, ok)

	tracePipeline, ok := pipelines[tracesPipelineRef.id()]
	test.NotNil(t, tracePipeline)
	test.True(t, ok)

	logsPipeline, ok := pipelines[component.NewID(component.DataTypeLogs)]
	test.Nil(t, logsPipeline)
	test.False(t, ok)
}

func Test_PipelineWithName(t *testing.T) {
	c := &Config{}
	hcpMetricsPipelineRef := c.NewPipelineWithName(component.DataTypeMetrics, "hcp")
	forwarderMetricsPipelineRef := c.NewPipelineWithName(component.DataTypeMetrics, "forwarder")
	metricsPipelineRef := c.NewPipeline(component.DataTypeMetrics)

	c.NewExporter(component.NewIDWithName("foo", "hcp"), hcpMetricsPipelineRef)
	c.NewExporter(component.NewID("foo"), forwarderMetricsPipelineRef)

	pipelines := c.Service.Pipelines

	test.NotEq(t, hcpMetricsPipelineRef, forwarderMetricsPipelineRef)

	metricPipeline, ok := pipelines[metricsPipelineRef.id()]
	test.NotNil(t, metricPipeline)
	test.True(t, ok)

	hcpPipeline, ok := pipelines[hcpMetricsPipelineRef.id()]
	test.True(t, ok)
	test.NotNil(t, hcpPipeline)
	test.NotNil(t, hcpPipeline.Exporters)
	test.Eq(t, hcpPipeline.Exporters[0], component.NewIDWithName("foo", "hcp"))

	forwarderPipeline, ok := pipelines[forwarderMetricsPipelineRef.id()]
	test.True(t, ok)
	test.NotNil(t, forwarderPipeline)
	test.NotNil(t, forwarderPipeline.Exporters)
	test.Eq(t, forwarderPipeline.Exporters[0], component.NewID("foo"))

	logsPipeline, ok := pipelines[component.NewID(component.DataTypeLogs)]
	test.Nil(t, logsPipeline)
	test.False(t, ok)
}
