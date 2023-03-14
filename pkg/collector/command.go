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

	Starts the telemetry-collector and runs until an interrupt is received. The
	collector can forward all metrics to an otlphttp endpoint or to the Hashicorp
	cloud platform.
`
)

// Command is the interface for running the collector
type Command struct {
	// ui is used for output. It should only be used if the logger has yet to
	// initialize. Otherwise always prefer the logger.
	ui cli.Ui

	// these are the initial flag values read in by the flag parser
	flagConfig *Config

	flags *flag.FlagSet
	help  string
}

// NewAgentCmd returns a new Agent command. It's mean to be only be called once
// to load configuration
func NewAgentCmd(ui cli.Ui) (*Command, error) {

	c := &Command{
		ui: ui,
	}

	c.flagConfig = &Config{Cloud: &Cloud{}}
	// Setup Flags
	c.flags = flag.NewFlagSet("", flag.ContinueOnError)
	c.flags.StringVar(&c.flagConfig.ConfigFile, COOConfigPathOpt, "", "Load configuration from a config file.")
	c.flags.StringVar(&c.flagConfig.Cloud.ClientID, HCPClientIDOpt, "", fmt.Sprintf("HCP Service Principal Client ID Environment variable %s", "HCP_CLIENT_ID"))
	c.flags.StringVar(&c.flagConfig.Cloud.ClientSecret, HCPClientSecretOpt, "", fmt.Sprintf("HCP Service Principal Client Secret Environment variable %s", "HCP_CLIENT_SECRET"))
	c.flags.StringVar(&c.flagConfig.Cloud.ResourceID, HCPResourceIDOpt, "", fmt.Sprintf("HCP Resource ID Environment variable %s", "HCP_RESOURCE_ID"))
	c.flags.StringVar(&c.flagConfig.HTTPCollectorEndpoint, COOtelHTTPEndpointOpt, "", fmt.Sprintf("OTLP HTTP endpoint to forward telemetry to Environment variable %s", "CO_OTEL_HTTP_ENDPOINT"))
	c.help = flags.Usage(help, c.flags)

	return c, nil
}

// loadConfiguration loads configuration in precedence order of highest first :
//
//  1. command line opts if specified
//  2. env variables if specified
//  3. file configuration if specified.
//
// This may or may not return a valid configuration. The function only parses the config from the values above.
// Validation is done in config.
func (c *Command) loadConfiguration(ctx context.Context, args []string, fileParser parser) (*Config, error) {
	logger := hclog.FromContext(ctx)
	logger.Debug("flag args passed to agent", "args", args)

	cfg := configFromEnvVars()

	// this parses the flags into c.flagConfig
	if err := c.flags.Parse(args); err != nil {
		logger.Debug("error parsing flags")
		return nil, err
	}

	// this controls for args like :
	//    $collector agent foobar -hcp-client-id
	//    $collector agent -hcp-client-id foobar
	remainingArgs := c.flags.Args()
	if len(remainingArgs) > 0 {
		return nil, fmt.Errorf("unexpected subcommands passed to command: %v", remainingArgs)
	}

	// c.flagConfig will override environment variable configuration
	if err := mergo.Merge(cfg, c.flagConfig, mergo.WithOverride); err != nil {
		return nil, err
	}

	if cfg.ConfigFile != "" {
		fileConfig, err := fileParser(cfg.ConfigFile)
		if err != nil {
			return nil, err
		}

		// file configuration will be overridden by the f+env variable config
		if err = mergo.Merge(fileConfig, cfg, mergo.WithOverride); err != nil {
			return nil, err
		}
		return fileConfig, nil
	}

	return cfg, nil
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
