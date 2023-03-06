package collector

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/go-hclog"
	"github.com/mitchellh/cli"
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

// Command is the interface for running the collector
type Command struct {
	// ui is used for output. It should only be used if the logger has yet to
	// initialize. Otherwise always prefer the logger.
	ui cli.Ui
}

// NewCollectorCmd returns a new Agent command
func NewCollectorCmd(ui cli.Ui) *Command {
	return &Command{
		ui: ui,
	}
}

// Synopsis gives details on how the collector runs
func (c *Command) Synopsis() string {
	return ""
}

// Help provides specifications on how to run the collector
func (c *Command) Help() string {
	return ""
}

// Run takes in args
func (c *Command) Run(args []string) int {

	ctx, cancel := context.WithCancel(context.Background())

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go handleSignal(sigCh, cancel)

	flags := flag.NewFlagSet("agent", flag.ContinueOnError)
	flags.Usage = func() { c.ui.Error(c.Help()) }

	logger := hclog.Default()
	var (
		configFile string

		// These are loaded from the environment or flags
		hcpClientID           string
		hcpClientSecret       string
		hcpResourceID         string
		httpCollectorEndpoint string
	)

	cfg := Config{
		HTTPCollectorEndpoint: httpCollectorEndpoint,
		ConfigFile:            configFile,
		Cloud: &Cloud{
			ClientSecret: hcpClientSecret,
			ClientID:     hcpClientID,
			ResourceID:   hcpResourceID,
		},
	}

	envVarString(HCPClientID, &hcpClientID)
	envVarString(HCPClientSecret, &hcpClientSecret)
	envVarString(HCPResourceID, &hcpResourceID)
	envVarString(COOtelHTTPEndpoint, &httpCollectorEndpoint)

	flags.StringVar(&configFile, "config-file", "", "Load configuration from a config file. Overrides environment and flag values")
	flags.StringVar(&hcpClientID, "hcp-client-id", "", fmt.Sprintf("HCP Service Principal Client ID \n\tEnvironment variable %s", "HCP_CLIENT_ID"))
	flags.StringVar(&hcpClientSecret, "hcp-client-secret", "", fmt.Sprintf("HCP Service Principal Client Secret \n\tEnvironment variable %s", "HCP_CLIENT_SECRET"))
	flags.StringVar(&hcpResourceID, "hcp-resource-id", "", fmt.Sprintf("HCP Resource ID \n\tEnvironment variable %s", "HCP_RESOURCE_ID"))
	flags.StringVar(&httpCollectorEndpoint, "http-collector-endpoint", "", fmt.Sprintf("OTLP HTTP endpoint to forward telemetry to \n\tEnvironment variable %s", "CO_OTEL_HTTP_ENDPOINT"))

	if err := flags.Parse(args); err != nil {
		cancel()
		logger.Error("error running collector", "error", err)
		c.ui.Error("error running collector")
		return -1
	}

	if err := loadConfig(cfg.ConfigFile, &cfg); err != nil {
		cancel()
		logger.Error("error loading configuration for collector", "error", err)
		c.ui.Error("error loading configuration for collector")
		return -1
	}

	if err := runSvc(ctx, cfg); err != nil {
		logger.Error("error running collector", "error", err)
		c.ui.Error("error running collector")
		return -1
	}

	return 0
}

func handleSignal(sigCh <-chan os.Signal, cancel func()) {
	sig := <-sigCh
	log.Println("Caught signal", sig.String(), ". Exiting")

	cancel()
}

func envVarString(envKey string, value *string) {
	if v, ok := os.LookupEnv(envKey); ok {
		*value = v
	}
}
