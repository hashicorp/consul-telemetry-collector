// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package metrics creates envoy grpc metricsserver. It will collect metrics,
// convert them from prometheus to OTLP. It will then push the OTLP metric onto the next component in an OTLP pipeline.
package metrics
