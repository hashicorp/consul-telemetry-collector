package confresolver

import (
	"testing"

	"github.com/shoenig/test"
	"go.opentelemetry.io/collector/component"
)

func Test_NewReceiver(t *testing.T) {
	cfg := &Config{}
	pipelineRef := cfg.NewPipeline(component.DataTypeTraces)
	ccfg := cfg.NewReceiver(pipelineRef, component.NewID("custom-trace-receiver"))
	otlp := cfg.NewReceiver(pipelineRef, component.NewID("otlp"))
	ccfg.Set("key", "value")
	ccfg.Set("endpoint", "http://localhost:8000")
	otlp.SetMap("protocols").SetMap("http")

	pipeline, ok := cfg.Service.Pipelines[pipelineRef.id()]
	test.True(t, ok)
	test.Len(t, 2, pipeline.Receivers)
	test.Eq(t, component.NewID("custom-trace-receiver"), pipeline.Receivers[0])
	test.Eq(t, component.NewID("otlp"), pipeline.Receivers[1])

	customReceiver, ok := cfg.Receivers[component.NewID("custom-trace-receiver")]
	test.True(t, ok)
	customComponentConfig, ok := customReceiver.(ComponentConfig)
	test.True(t, ok)
	test.Eq(t, ccfg, customComponentConfig)

	otlpReceiver, ok := cfg.Receivers[component.NewID("otlp")]
	test.True(t, ok)
	otlpComponentConfig, ok := otlpReceiver.(ComponentConfig)
	test.True(t, ok)
	test.Eq(t, otlp, otlpComponentConfig)
}
