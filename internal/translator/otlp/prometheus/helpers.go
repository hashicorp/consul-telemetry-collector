// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package prometheus

import (
	"strings"
	"time"

	"go.opentelemetry.io/collector/pdata/pcommon"
)

const (
	suffixCount   = "_count"
	suffixBucket  = "_bucket"
	suffixSum     = "_sum"
	suffixTotal   = "_total"
	suffixInfo    = "_info"
	suffixCreated = "_created"
)

var (
	suffixes = []string{suffixCreated, suffixBucket, suffixInfo, suffixSum, suffixCount}
)

func timestampFromMs(timestampMs int64) pcommon.Timestamp {
	t := time.Unix(0, timestampMs*int64(time.Millisecond))
	return pcommon.NewTimestampFromTime(t)
}

func normalizeName(name string) string {
	for _, suffix := range suffixes {
		if strings.HasSuffix(name, suffix) && name != suffix {
			return strings.TrimSuffix(name, suffix)
		}
	}
	return name
}
