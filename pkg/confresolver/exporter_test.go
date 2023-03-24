package confresolver

import (
	"testing"

	"github.com/shoenig/test"
	"go.opentelemetry.io/collector/component"
)

func Test_NewExporter(t *testing.T) {
	cfg := &Config{}
	pipelineRef := cfg.NewPipeline(component.DataTypeTraces)
	ccfg := cfg.NewExporter(component.NewID("custom-trace-exporter"), pipelineRef)
	otlp := cfg.NewExporter(component.NewID("otlp"), pipelineRef)
	ccfg.Set("key", "value")
	ccfg.Set("endpoint", "http://localhost:8000")
	otlp.SetMap("protocols").SetMap("http")

	pipeline, ok := cfg.Service.Pipelines[pipelineRef.id()]
	test.True(t, ok)
	test.Len(t, 2, pipeline.Exporters)
	test.Eq(t, component.NewID("custom-trace-exporter"), pipeline.Exporters[0])
	test.Eq(t, component.NewID("otlp"), pipeline.Exporters[1])

	customExporter, ok := cfg.Exporters[component.NewID("custom-trace-exporter")]
	test.True(t, ok)
	customComponentConfig, ok := customExporter.(ComponentConfig)
	test.True(t, ok)
	test.Eq(t, ccfg, customComponentConfig)

	otlpReceiver, ok := cfg.Exporters[component.NewID("otlp")]
	test.True(t, ok)
	otlpComponentConfig, ok := otlpReceiver.(ComponentConfig)
	test.True(t, ok)
	test.Eq(t, otlp, otlpComponentConfig)
}
