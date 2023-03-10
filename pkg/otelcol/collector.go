package otelcol

import (
	"context"

	"go.opentelemetry.io/collector/otelcol"
)

// Collector is an interface that is satisfied by the otelcol.Collector struct.
// This allows us to wrap the opentelemetry collector and not necessarily run it ourselves
type Collector interface {
	Run(context.Context) error
	GetState() otelcol.State
	Shutdown()
}
