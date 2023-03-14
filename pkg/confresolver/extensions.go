package confresolver

import "go.opentelemetry.io/collector/component"

// NewExtensions adds a new extension component to the overall configuration and
// adds the component to the service extensions. It returns the component's
// configuration for further configuring.
func (c *Config) NewExtensions(id component.ID) ComponentConfig {
	var ccfg componentConfig

	if c.Extensions == nil {
		c.Extensions = make(telemetryComponents)
	}

	c.Service.Extensions, ccfg = addComponent(c.Extensions, id, c.Service.Extensions)
	return ccfg
}
