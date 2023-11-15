// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package providers

import "time"

// SharedParams holds shared configuration parameters
type SharedParams struct {
	EnvoyPort    int
	MetricsPort  int
	BatchTimeout time.Duration
}
