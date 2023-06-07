// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package otel handles the configuration and lifecycle of the opentelemetry-collector
// Its' purpose is to generate a resolver setting that incorporates multiple providers
// collector consumes otel/providers which consume otel/config
package otel
