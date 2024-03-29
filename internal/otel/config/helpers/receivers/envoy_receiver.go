// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package receivers holds the type of receivers that consul telemetery supports
package receivers

import (
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configgrpc"

	"github.com/hashicorp/consul-telemetry-collector/receivers/envoyreceiver"
)

// EnvoyReceiverID is the component id of the otlp receiver.
var EnvoyReceiverID component.ID = component.NewID(envoyreceiver.ID)

// Config is the configuration for the supported protocols.
type Config struct {
	GRPC *configgrpc.GRPCServerSettings `mapstructure:"grpc,omitempty"`
}

// EnvoyReceiverCfg  generates the config for an otlp receiver.
func EnvoyReceiverCfg(listenerPort int) *envoyreceiver.Config {
	defaults := envoyreceiver.NewFactory().CreateDefaultConfig().(*envoyreceiver.Config)
	defaults.GRPC.NetAddr.Endpoint = fmt.Sprintf("127.0.0.1:%d", listenerPort)

	return defaults
}
