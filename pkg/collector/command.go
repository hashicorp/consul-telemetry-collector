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
	"github.com/imdario/mergo"
	"github.com/mitchellh/cli"

	"github.com/hashicorp/consul-telemetry-collector/pkg/flags"
)

const (
	synopsis = "Runs telemetry collector in agent mode"
	help     = `
Usage: consul-telemetry-collector agent [options]

	Starts the Consul agent and runs until an interrupt is received. The
	agent represents a single node in a cluster.
`
)

// Command is the interface for running the collector
type Command struct {
	// ui is used for output. It should only be used if the logger has yet to
	// initialize. Otherwise always prefer the logger.
	ui cli.Ui

	flags *flag.FlagSet
	help  string
	cfg   *Config
}

// NewAgentCmd returns a new Agent command
func NewAgentCmd(ui cli.Ui) (*Command, error) {
	c := &Command{
		ui: ui,
	}

	var (
		configFilePath, hcpClientID, hcpClientSecret, hcpResourceID, httpCollectorEndpoint string
	)

	// Setup Flags
	c.flags = flag.NewFlagSet("", flag.ContinueOnError)

	c.flags.StringVar(&configFilePath, COOConfigPathOpt, "", "Load configuration from a config file.")
	c.flags.StringVar(&hcpClientID, HCPClientIDOpt, "", fmt.Sprintf("HCP Service Principal Client ID Environment variable %s", "HCP_CLIENT_ID"))
	c.flags.StringVar(&hcpClientSecret, HCPClientSecretOpt, "", fmt.Sprintf("HCP Service Principal Client Secret Environment variable %s", "HCP_CLIENT_SECRET"))
	c.flags.StringVar(&hcpResourceID, HCPResourceIDOpt, "", fmt.Sprintf("HCP Resource ID Environment variable %s", "HCP_RESOURCE_ID"))
	c.flags.StringVar(&httpCollectorEndpoint, COOtelHTTPEndpointOpt, "", fmt.Sprintf("OTLP HTTP endpoint to forward telemetry to Environment variable %s", "CO_OTEL_HTTP_ENDPOINT"))

	defaultToEnv(&configFilePath, COOConfigPath)
	defaultToEnv(&hcpClientID, HCPClientID)
	defaultToEnv(&hcpClientSecret, HCPClientSecret)
	defaultToEnv(&hcpResourceID, HCPResourceID)
	defaultToEnv(&httpCollectorEndpoint, COOtelHTTPEndpoint)

	c.cfg = &Config{
		HTTPCollectorEndpoint: &httpCollectorEndpoint,
		ConfigFile:            &configFilePath,
		Cloud: &Cloud{
			ClientSecret: &hcpClientSecret,
			ClientID:     &hcpClientID,
			ResourceID:   &hcpResourceID,
		},
	}
	c.help = flags.Usage(help, c.flags)

	return c, nil
}

func defaultToEnv(param *string, envKey string) {
	if v, ok := os.LookupEnv(envKey); ok {
		*param = v
	}
}

// loadConfiguration loads configuration in precdence order of 1. CLI flag options 2. OS env variables 3. file configuration
// This may or may not return a valid configuration. The function only parses the config from the values above.
// Validation is done in config.
func (c *Command) loadConfiguration(ctx context.Context, args []string, fileParser parser) (*Config, error) {
	logger := hclog.FromContext(ctx)
	logger.Debug("cli args passed to agent", "args", args)

	if err := c.flags.Parse(args); err != nil {
		logger.Debug("error parsing flags")
		return nil, err
	}

	// this controls for args like :
	// 		$collector agent foobar -hcp-client-id
	//		$collector agent -hcp-client-id foobar
	remainingArgs := c.flags.Args()
	if len(remainingArgs) > 0 {
		return nil, fmt.Errorf("unexpected subcommands passed to command: %v", remainingArgs)
	}

	if *c.cfg.ConfigFile != "" {
		fileConfig, err := fileParser(c.cfg.ConfigFile)
		if err != nil {
			return nil, err
		}

		// Precedence
		// 1. command line opts if specified
		// 2. env variables if specified
		// 3. file config
		if err = mergo.Merge(fileConfig, c.cfg, mergo.WithOverride); err != nil {
			return nil, err
		}
		return fileConfig, nil
	}

	return c.cfg, nil
}

// Synopsis gives details on how the collector runs
func (c *Command) Synopsis() string {
	return synopsis
}

// Help provides specifications on how to run the collector
func (c *Command) Help() string {
	return c.help
}

// Run takes in args and runs the collector agent as an OTEL collector
func (c *Command) Run(args []string) int {
	logger := hclog.Default().Named("consul-collector")
	ctx, cancel := context.WithCancel(hclog.WithContext(context.Background(), logger))
	logger.Info("debugging args", "args", args)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go handleSignal(sigCh, cancel)

	// load the configuration
	cfg, err := c.loadConfiguration(ctx, args, parseFile)
	if err != nil {
		cancel()
		logger.Error("error loading configuration", "error", err)
		return -1
	}

	// validate the loaded configuration is valid
	if err := cfg.validate(); err != nil {
		logger.Error("configuration is invalid", "error", err)
		return -1
	}

	// run the service
	if err := runSvc(ctx, cfg); err != nil {
		cancel()
		logger.Error("error running collector", "error", err)
		return -1
	}

	return 0
}

func handleSignal(sigCh <-chan os.Signal, cancel func()) {
	sig := <-sigCh
	log.Println("Caught signal", sig.String(), ". Exiting")
	cancel()
}
