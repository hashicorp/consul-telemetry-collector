package main

import (
	"fmt"

	"github.com/hashicorp/consul-telemetry-collector/pkg/version"
)

func main() {
	fmt.Printf("Consul Telemetry Collector v%s\n", version.GetHumanVersion())
}
