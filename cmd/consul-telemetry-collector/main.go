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
	// HCPClientID is the environment variable for the hcpClientID
	HCPClientID = "HCP_CLIENT_ID"
	// HCPClientSecret is the environment variable for the hcpClientSecret
	HCPClientSecret = "HCP_CLIENT_SECRET"
	// HCPResourceID is the environment variable for the hcpResourceID
	HCPResourceID = "HCP_RESOURCE_ID"
	// COOtelHTTPEndpoint is the environment variable for the OpenTelemetry HTTP Endpoint we forward metrics metrics to
	COOtelHTTPEndpoint = "CO_OTEL_HTTP_ENDPOINT"
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

	envVarString(HCPClientID, &hcpClientID)
	envVarString(HCPClientSecret, &hcpClientSecret)
	envVarString(HCPResourceID, &hcpResourceID)
	envVarString(COOtelHTTPEndpoint, &httpCollectorEndpoint)

	flag.BoolVar(&printVersion, "version", false, "Print the build version and exit")
	flag.StringVar(&configFile, "config-file", "", "Load configuration from a config file. Overrides environment and flag values")
	flag.StringVar(&hcpClientID, "hcp-client-id", "", fmt.Sprintf("HCP Service Principal Client ID \n\tEnvironment variable %s", "HCP_CLIENT_ID"))
	flag.StringVar(&hcpClientSecret, "hcp-client-secret", "", fmt.Sprintf("HCP Service Principal Client Secret \n\tEnvironment variable %s", "HCP_CLIENT_SECRET"))
	flag.StringVar(&hcpResourceID, "hcp-resource-id", "", fmt.Sprintf("HCP Resource ID \n\tEnvironment variable %s", "HCP_RESOURCE_ID"))
	flag.StringVar(&httpCollectorEndpoint, "http-collector-endpoint", "", fmt.Sprintf("OTLP HTTP endpoint to forward telemetry to \n\tEnvironment variable %s", "CO_OTEL_HTTP_ENDPOINT"))

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

	go handleSignal(sigCh, cancel)

	cfg := flags()

	if printVersion {
		fmt.Printf("Consul Telemetry Collector v%s\n", version.GetHumanVersion())
		return
	}

	if err := collector.Run(ctx, cfg); err != nil {
		log.Fatal(err)
	}
}

func handleSignal(sigCh <-chan os.Signal, cancel func()) {
	sig := <-sigCh
	log.Println("Caught signal", sig.String(), ". Exiting")
	cancel()
}
