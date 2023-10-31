package providers

import "time"

type SharedParams struct {
	EnvoyPort    int
	MetricsPort  int
	BatchTimeout time.Duration
}
