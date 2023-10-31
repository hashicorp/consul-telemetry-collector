package providers

import "time"

// SharedParams holds shared configuration parameters
type SharedParams struct {
	EnvoyPort    int
	MetricsPort  int
	BatchTimeout time.Duration
}
