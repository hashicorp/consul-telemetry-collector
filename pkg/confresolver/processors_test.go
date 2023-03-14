package confresolver

import (
	"testing"

	"github.com/shoenig/test"
	"github.com/shoenig/test/must"
	"go.opentelemetry.io/collector/component"
)

func Test_NewProcessor(t *testing.T) {
	cfg := &Config{}
	pipelineRef := cfg.NewPipeline(component.DataTypeMetrics)
	ccfg := cfg.NewProcessor(pipelineRef, component.NewID("batch"))

	ccfg.Set("send_batch_size", 10000)
	ccfg.Set("timeout", "10s")

	pipeline, ok := cfg.Service.Pipelines[pipelineRef.id()]
	test.True(t, ok)
	test.Len(t, 1, pipeline.Processors)
	test.Eq(t, component.NewID("batch"), pipeline.Processors[0])

	processor := cfg.Processors[component.NewID("batch")]
	processorCfg, ok := processor.(componentConfig)
	must.True(t, ok)

	componentConfig, ok := ccfg.(componentConfig)
	must.True(t, ok)

	test.Eq(t, componentConfig["send_batch_size"], processorCfg["send_batch_size"])
	test.Eq(t, componentConfig["timeout"], processorCfg["timeout"])
}
