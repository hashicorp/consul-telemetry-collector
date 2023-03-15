package collector

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/shoenig/test"
)

func setupEnv(t *testing.T, env map[string]string) {
	t.Helper()
	for k, v := range env {
		os.Setenv(k, v)
	}
	t.Cleanup(func() {
		for k := range env {
			os.Unsetenv(k)
		}
	})
}

func testConfig() *Config {
	return &Config{
		Cloud: &Cloud{},
	}
}

func wrapOpt(s string) string {
	return fmt.Sprintf("-%s", s)
}

func Test_loadConfiguration(t *testing.T) {

	for name, tc := range map[string]struct {
		configPath string
		// this is the file config returned from file parsing.
		mutateFileConfig func(*Config)
		// args are the flags passed in by the user
		args []string
		// env is the os environment variables set by the user
		env map[string]string

		// the is the expected config
		mutateExpected func(*Config)

		err error
	}{
		"Invalidflags": {
			args: []string{"-hcp-client-id 123"},
			err:  errors.New("flag provided but not defined: -hcp-client-id 123"),
		},
		"ExtraSubCommands": {
			args: []string{"foobar"},
			err:  errors.New("unexpected subcommands passed to command: [foobar]"),
		},
		"ExtraSubCommandsDisordered": {
			args: []string{"-hcp-client-id", "123", "foobar"},
			err:  errors.New("unexpected subcommands passed to command: [foobar]"),
		},
		"ErrorNeverReachedHelp": {
			args: []string{"-h"},
			err:  errors.New("flag: help requested"),
		},
		// Note that this isn't a valid set of config but we will split parsing configuration
		// from validating config so it is considered a success to load config with 0 flag opts, 0 env var
		// 0 file paths.
		"Success": {},
		"SuccessWithAllEnv": {
			env: map[string]string{
				HCPClientID:        "id",
				HCPClientSecret:    "sec",
				HCPResourceURL:     "rid",
				COOtelHTTPEndpoint: "ep",
				COOConfigPath:      "fp",
			},
			mutateExpected: func(c *Config) {
				c.Cloud.ClientID = "id"
				c.Cloud.ClientSecret = "sec"
				c.Cloud.ResourceURL = "rid"
				c.HTTPCollectorEndpoint = "ep"
				c.ConfigFile = "fp"
			},
		},
		"SuccessWithCliOptsPrecedenceOverEnvVariables": {
			args: []string{
				wrapOpt(HCPClientIDOpt),
				"flagid",
				wrapOpt(HCPClientSecretOpt),
				"flagsec",
				wrapOpt(HCPResourceURLOpt),
				"flagrid",
				wrapOpt(COOtelHTTPEndpointOpt),
				"flagep",
				wrapOpt(COOConfigPathOpt),
				"flagfp",
			},
			env: map[string]string{
				HCPClientID:        "id",
				HCPClientSecret:    "sec",
				HCPResourceURL:     "rid",
				COOtelHTTPEndpoint: "ep",
				COOConfigPath:      "fp",
			},
			mutateExpected: func(c *Config) {
				c.Cloud.ClientID = "flagid"
				c.Cloud.ClientSecret = "flagsec"
				c.Cloud.ResourceURL = "flagrid"
				c.HTTPCollectorEndpoint = "flagep"
				c.ConfigFile = "flagfp"
			},
		},
		"SuccessWithEnvVariablePrecedenceOverFileCfg": {
			env: map[string]string{
				HCPClientID:        "id",
				HCPClientSecret:    "sec",
				HCPResourceURL:     "rid",
				COOtelHTTPEndpoint: "ep",
				COOConfigPath:      "fp",
			},
			mutateFileConfig: func(c *Config) {
				c.Cloud.ClientID = "fid"
				c.Cloud.ClientSecret = "fsec"
				c.Cloud.ResourceURL = "fid"
				c.HTTPCollectorEndpoint = "fep"
				c.ConfigFile = "fp"
			},
			mutateExpected: func(c *Config) {
				c.Cloud.ClientID = "id"
				c.Cloud.ClientSecret = "sec"
				c.Cloud.ResourceURL = "rid"
				c.HTTPCollectorEndpoint = "ep"
				c.ConfigFile = "fp"
			},
		},
		"SuccessWithCliOptsPrecedenceOverEnvVariablesOverFileCfg": {
			args: []string{
				wrapOpt(HCPClientIDOpt),
				"flagid",
				wrapOpt(HCPClientSecretOpt),
				"flagsec",
				wrapOpt(HCPResourceURLOpt),
				"flagrid",
				wrapOpt(COOtelHTTPEndpointOpt),
				"flagep",
				wrapOpt(COOConfigPathOpt),
				"flagfp",
			},
			env: map[string]string{
				HCPClientID:        "id",
				HCPClientSecret:    "sec",
				HCPResourceURL:     "rid",
				COOtelHTTPEndpoint: "ep",
				COOConfigPath:      "fp",
			},
			mutateFileConfig: func(c *Config) {
				c.Cloud.ClientID = "fid"
				c.Cloud.ClientSecret = "fsec"
				c.Cloud.ResourceURL = "frid"
				c.HTTPCollectorEndpoint = "fep"
				c.ConfigFile = "ffp"
			},
			mutateExpected: func(c *Config) {
				c.Cloud.ClientID = "flagid"
				c.Cloud.ClientSecret = "flagsec"
				c.Cloud.ResourceURL = "flagrid"
				c.HTTPCollectorEndpoint = "flagep"
				c.ConfigFile = "flagfp"
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			setupEnv(t, tc.env)
			c, err := NewAgentCmd(&cli.BasicUi{})
			test.NoError(t, err)
			fileConfig := testConfig()
			if tc.mutateFileConfig != nil {
				tc.mutateFileConfig(fileConfig)
			}
			parser := func(string) (*Config, error) {
				return fileConfig, tc.err
			}
			args := []string{}
			if tc.args != nil {
				args = tc.args
			}
			config, err := c.loadConfiguration(context.Background(), args, parser)
			if tc.err != nil {
				test.Error(t, err)
				test.ErrorContains(t, err, tc.err.Error())
				return
			}
			test.NoError(t, err)
			expectedCfg := testConfig()
			if tc.mutateExpected != nil {
				tc.mutateExpected(expectedCfg)
			}
			test.Eq(t, config, expectedCfg)
		})
	}
}
