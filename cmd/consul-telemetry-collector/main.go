package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/consul-telemetry-collector/pkg/collector"
	"github.com/hashicorp/consul-telemetry-collector/pkg/version"
)

var (
	printVersion          bool
	hcpClientID           string
	hcpClientSecret       string
	hcpResourceID         string
	httpCollectorEndpoint string
)

func flags() {
	// We'll want to implement the flag.Var interface which we can use to figure out if these values actually got set,
	// fallback to environment variables, and override configuration files (only _really_ necessary for the cloud secrets)
	flag.BoolVar(&printVersion, "version", false, "")
	flag.StringVar(&hcpClientID, "hcp-client-id", "", "")
	flag.StringVar(&hcpClientSecret, "hcp-client-secret", "", "")
	flag.StringVar(&hcpResourceID, "hcp-resource-id", "", "")
	flag.StringVar(&httpCollectorEndpoint, "http-collector-endpoint", "", "")
}

func main() {
	flags()

	if printVersion {
		fmt.Printf("Consul Telemetry Collector v%s\n", version.GetHumanVersion())
		return
	}

	cfg := collector.Config{
		HTTPCollectorEndpoint: httpCollectorEndpoint,
		Cloud: collector.Cloud{
			ClientSecret: hcpClientSecret,
			ClientID:     hcpClientID,
		},
	}

	ctx, cancel := context.WithCancel(context.Background())

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		log.Println("Caught", sig.String(), ". exiting")
		cancel()
	}()

	if err := collector.Run(ctx, cfg); err != nil {
		log.Fatal(err)
	}

}
