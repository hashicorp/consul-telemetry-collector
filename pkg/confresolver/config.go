package confresolver

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/service"
)

// Config is a helper type to create a new opentelemetry server configuration
type Config struct {
	Receivers  components     `mapstructure:"receivers"`
	Exporters  components     `mapstructure:"exporters"`
	Processors components     `mapstructure:"processors"`
	Connectors components     `mapstructure:"connectors"`
	Extensions components     `mapstructure:"extensions"`
	Service    service.Config `mapstructure:"service"`
}

type components map[component.ID]interface{}
type componentConfig map[string]interface{}

// ComponentConfig is an interface that lets us set key/value entries or child maps on the component
type ComponentConfig interface {
	Set(k, v string)
	Map(k string) ComponentConfig
}

func (t componentConfig) Set(k, v string) {
	t[k] = v
}

func (t componentConfig) Map(k string) ComponentConfig {
	tc := make(componentConfig)
	t[k] = tc
	return tc
}

func addComponent(pipelineComponentConfig map[component.ID]interface{}, id component.ID,
	pipelineComponent []component.ID) ([]component.ID, componentConfig) {
	// create the new config
	ccfg := make(componentConfig)
	pipelineComponentConfig[id] = ccfg

	// add to the pipeline slice (and make sure that's not nil)
	if pipelineComponent == nil {
		pipelineComponent = make([]component.ID, 0, 1)
	}
	pipelineComponent = append(pipelineComponent, id)

	// return both the slice and config
	return pipelineComponent, ccfg
}