package extensions

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/extension/ballastextension"
)

const (
	ballastName = "memory_ballast"
)

// BallastID is the component id of the ballast extension
var BallastID component.ID = component.NewID(ballastName)

// BallastCfg  generates the config for a ballast config
func BallastCfg() (component.ID, *ballastextension.Config) {
	return BallastID, &ballastextension.Config{
		SizeInPercentage: 10,
	}
}
