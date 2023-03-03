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

const (
	HCP_CLIENT_ID         = "HCP_CLIENT_ID"
	HCP_CLIENT_SECRET     = "HCP_CLIENT_SECRET"
	HCP_RESOURCE_ID       = "HCP_RESOURCE_ID"
	CO_OTEL_HTTP_ENDPOINT = "CO_OTEL_HTTP_ENDPOINT"
)

var (
	// If this flag is set, print the human read-able collector version and exit
	printVersion bool
)

// flags will set CLI flags and environment variables into the config for later use and to potentially override with
// a configuration file
func flags() collector.Config {
	var (
		configFile string

		// These are loaded from the environment or flags
		hcpClientID           string
		hcpClientSecret       string
		hcpResourceID         string
		httpCollectorEndpoint string
	)

	cfg := collector.Config{
		HTTPCollectorEndpoint: httpCollectorEndpoint,
		ConfigFile:            configFile,
		Cloud: &collector.Cloud{
			ClientSecret: hcpClientSecret,
			ClientID:     hcpClientID,
			ResourceID:   hcpResourceID,
		},
	}

	envVarString(HCP_CLIENT_ID, &hcpClientID)
	envVarString(HCP_CLIENT_SECRET, &hcpClientSecret)
	envVarString(HCP_RESOURCE_ID, &hcpResourceID)
	envVarString(CO_OTEL_HTTP_ENDPOINT, &httpCollectorEndpoint)

	flag.BoolVar(&printVersion, "version", false, "Print the build version and exit")
	flag.StringVar(&configFile, "config-file", "", "Load configuration from a config file. Overrides environment and flag values")
	flag.StringVar(&hcpClientID, "hcp-client-id", "", fmt.Sprintf("HCP Service Principal Client ID \n\tEnvironment variable %s", "HCP_CLIENT_ID"))
	flag.StringVar(&hcpClientSecret, "hcp-client-secret", "", fmt.Sprintf("HCP Service Principal Client Secret \n\tEnvironment variable %s", "HCP_CLIENT_SECRET"))
	flag.StringVar(&hcpResourceID, "hcp-resource-id", "", fmt.Sprintf("HCP Resource ID \n\tEnvironment variable %s", "HCP_RESOURCE_ID"))
	flag.StringVar(&httpCollectorEndpoint, "http-collector-endpoint", "", fmt.Sprintf("OTLP HTTP endpoint to forward telemetry to \n\tEnvironment variable %s", "CO_OTEL_HTTP_ENDPOINT"))

	// flags will override environment variables set in environmentConfig
	flag.Parse()

	return cfg
}

func envVarString(envKey string, value *string) {
	if v, ok := os.LookupEnv(envKey); ok {
		*value = v
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		log.Println("Caught", sig.String(), ". exiting")
		cancel()
	}()

	cfg := flags()

	if printVersion {
		fmt.Printf("Consul Telemetry Collector v%s\n", version.GetHumanVersion())
		return
	}

	if err := collector.Run(ctx, cfg); err != nil {
		log.Fatal(err)
	}
}
