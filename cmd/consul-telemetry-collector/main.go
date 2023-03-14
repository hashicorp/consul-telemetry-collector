package main

import (
	"os"

	"github.com/hashicorp/consul-telemetry-collector/pkg/collector"
	"github.com/hashicorp/consul-telemetry-collector/pkg/version"
	"github.com/hashicorp/go-hclog"
	"github.com/mitchellh/cli"
)

const (
	appName = "consul-collector"
)

func main() {
	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	commands := map[string]cli.CommandFactory{
		"agent": func() (cli.Command, error) {
			return collector.NewAgentCmd(ui)
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
