// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package agent

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"go.uber.org/multierr"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/hcl/v2/hclsimple"
)

var (
	errNoConfigurationProvided = errors.New("no configuration provided: see usage")
	errCloudConfigInvalid      = errors.New("cloud configuration is not valid")
)

func configFromEnvVars() *Config {
	return &Config{
		Cloud: &Cloud{
			ClientID:     os.Getenv(HCPClientID),
			ClientSecret: os.Getenv(HCPClientSecret),
			ResourceID:   os.Getenv(HCPResourceID),
		},
		ConfigFile:            os.Getenv(COOConfigPath),
		HTTPCollectorEndpoint: os.Getenv(COOtelHTTPEndpoint),
	}
}

// used to parse a file path and return a configuration.
type parser func(string) (*Config, error)

// ParseFile parses the given file for a configuration. The file is expected
// to be in JSON or HCL format.
func parseFile(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	// wrap our file in a 1mb limit reader
	mb := int64(1024 * 1024 * 1024)
	r := io.LimitReader(f, mb)
	cfg, err := readConfiguration(r, filename)
	cerr := f.Close()
	if err != nil {
		return nil, multierr.Append(err, cerr)
	}
	return cfg, cerr
}

func readConfiguration(reader io.Reader, filename string) (*Config, error) {
	cfg := &Config{}
	buffer, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}
	// decode needs filename for parsing and bytes passed to it.
	err = hclsimple.Decode(filename, buffer, nil, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed parsing config file: %w", err)
	}

	if err = parseExportConfig(cfg); err != nil {
		return nil, fmt.Errorf("failed parsing config file: %w", err)
	}

	return cfg, nil
}

func parseExportConfig(c *Config) error {
	if c.ExporterConfig != nil {
		if c.ExporterConfig.Timeout != "" {
			d, err := time.ParseDuration(c.ExporterConfig.Timeout)
			if err != nil {
				return fmt.Errorf("failed to parse export config. unable to parse timeout %s %w", c.ExporterConfig.Timeout, err)
			}
			c.ExporterConfig.timeoutDuration = d
		}
	}
	return nil
}

// Config is the global collector configuration.
type Config struct {
	Cloud                 *Cloud `hcl:"cloud,block"`
	HTTPCollectorEndpoint string `hcl:"http_collector_endpoint,optional"`
	ConfigFile            string
	ExporterConfig        *ExporterConfig `hcl:"exporter_config,block"`
}

// Cloud is the HCP Cloud configuration.
type Cloud struct {
	ClientID     string `hcl:"client_id,optional"`
	ClientSecret string `hcl:"client_secret,optional"`
	ResourceID   string `hcl:"resource_id,optional"`
}

// ExporterConfig holds
type ExporterConfig struct {
	Type            string            `hcl:"type,label"`
	Headers         map[string]string `hcl:"headers,optional"`
	Endpoint        string            `hcl:"endpoint"`
	Timeout         string            `hcl:"timeout,optional"`
	timeoutDuration time.Duration
}

// IsEnabled checks if the Cloud config is enabled. It returns false if the ClientID,
// ClientSecret and ResourceID are all empty.
func (c *Cloud) IsEnabled() bool {
	if c == nil {
		return false
	}

	if c.ClientSecret != "" || c.ClientID != "" || c.ResourceID != "" {
		return true
	}

	return false
}

// validate that, if the Cloud config is enabled, all required fields are set.
// Return an error describing missing fields.
func (c *Cloud) validate() error {
	if c == nil {
		return nil
	}

	if !c.IsEnabled() {
		return nil
	}

	missing := []string{}
	if c.ClientID == "" {
		missing = append(missing, "client_id")
	}
	if c.ClientSecret == "" {
		missing = append(missing, "client_secret")
	}
	if c.ResourceID == "" {
		missing = append(missing, "resource_id")
	}
	if len(missing) > 0 {
		return fmt.Errorf("%w: missing %s", errCloudConfigInvalid, strings.Join(missing, ", "))
	}

	return nil
}

func (c *Config) validate() error {
	if c == nil {
		return errNoConfigurationProvided
	}

	if c.Cloud == nil {
		return nil
	}

	return c.Cloud.validate()
}

func (c *Config) logDeprecations(logger hclog.Logger) {
	const deprecatedWarning = "'%s' is deprecated and will be removed in a future release. Use '%s' instead."
	const conflictingConfig = "deprecated field '%s' and supported config '%s' are both configured. Using '%s'"
	if c.HTTPCollectorEndpoint != "" {
		logger.Warn(deprecatedWarning, "http_collector_endpoint", "exporter_config")
	}

	if c.ExporterConfig != nil && c.HTTPCollectorEndpoint != "" {
		logger.Warn(conflictingConfig, "http_collector_endpoint", "exporter_config", "exporter_config")
	}
}
