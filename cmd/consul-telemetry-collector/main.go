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
	// static loader
	printVersion bool
	configFile   string

	// These are loaded from the environment or flag
	hcpClientID           string
	hcpClientSecret       string
	hcpResourceID         string
	httpCollectorEndpoint string
)

func flags() {
	flag.BoolVar(&printVersion, "version", false, "Print the build version and exit")
	flag.StringVar(&configFile, "config-file", "", "Load configuration from a config file. "+
		"Overrides environment and flag values")

	StringVar(&hcpClientID, "hcp-client-id", "", "HCP Service Principal Client ID", "HCP_CLIENT_ID")
	StringVar(&hcpClientSecret, "hcp-client-secret", "", "HCP Service Principal Client Secret", "HCP_CLIENT_SECRET")
	StringVar(&hcpResourceID, "hcp-resource-id", "", "HCP Resource ID", "HCP_RESOURCE_ID")
	StringVar(&httpCollectorEndpoint, "http-collector-endpoint", "", "OTLP HTTP endpoint to forward telemetry to",
		"CO_OTEL_HTTP_ENDPOINT")

}

func main() {
	flags()

	if printVersion {
		fmt.Printf("Consul Telemetry Collector v%s\n", version.GetHumanVersion())
		return
	}

	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
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

func loadConfig() (collector.Config, error) {
	cfg := collector.Config{
		HTTPCollectorEndpoint: httpCollectorEndpoint,
		ConfigFile:            configFile,
		Cloud: &collector.Cloud{
			ClientSecret: hcpClientSecret,
			ClientID:     hcpClientID,
			ResourceID:   hcpResourceID,
		},
	}

	if configFile != "" {
		err := collector.LoadConfig(configFile, &cfg)
		return cfg, err
	}

	return cfg, nil
}
