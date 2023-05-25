package main

import (
	"os"

	"github.com/mitchellh/cli"

	"github.com/hashicorp/consul-telemetry-collector/internal/agent"
	"github.com/hashicorp/consul-telemetry-collector/internal/version"
	"github.com/hashicorp/go-hclog"
)

const (
	appName = "consul-telemetry-collector"
)

func main() {
	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	commands := map[string]cli.CommandFactory{
		"agent": func() (cli.Command, error) {
			return agent.NewAgentCmd(ui)
		},
	}

	// Build and run the CLI
	cli := &cli.CLI{
		Name:                       appName,
		Version:                    version.GetHumanVersion(),
		Args:                       os.Args[1:],
		Commands:                   commands,
		Autocomplete:               true,
		AutocompleteNoDefaultFlags: true,
		HelpFunc:                   cli.BasicHelpFunc(appName),
		HelpWriter:                 os.Stdout,
	}

	exitStatus, err := cli.Run()
	if err != nil {
		hclog.Default().Error("failed to run CLI", "error", err)
	}

	os.Exit(exitStatus)
}
