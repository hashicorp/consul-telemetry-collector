package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	telemetrycollector "github.com/hashicorp/consul-telemetry-collector/internal/telemetry-collector"
	"github.com/hashicorp/consul-telemetry-collector/pkg/version"
)

var (
	printVersion bool
)

func flags() {
	flag.BoolVar(&printVersion, "version", false, "")
	flag.Parse()
}

func main() {
	flags()

	if printVersion {
		fmt.Printf("Consul Telemetry Collector v%s\n", version.GetHumanVersion())
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		log.Println("Caught", sig.String())
		cancel()
	}()

	err := telemetrycollector.Run(ctx)
	if err != nil {
		cancel()
		log.Fatal(err)
	}
}
