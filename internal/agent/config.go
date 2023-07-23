// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package agent

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"go.uber.org/multierr"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

var (
	errNoConfigurationProvided = errors.New("no configuration provided: see usage")
	errCloudConfigInvalid      = errors.New("cloud configuration is not valid")
)

type errCloudConfigInvalidField struct {
	cloudCfg Cloud
}

// Error returns a string describing the invalid fields in the Cloud config. Each invalid field is separated by a comma
func (e *errCloudConfigInvalidField) Error() string {
	sb := strings.Builder{}
	if e.cloudCfg.ClientID == "" {
		sb.WriteString("client_id")
	}
	if e.cloudCfg.ClientSecret == "" {
		if sb.Len() > 0 {
			sb.WriteString(",")
		}
		sb.WriteString("client_secret,")
	}
	if e.cloudCfg.ResourceID == "" {
		if sb.Len() > 0 {
			sb.WriteString(",")
		}
		sb.WriteString("resource_id")
	}
	return fmt.Sprintf("%w: missing %s", errCloudConfigInvalid, sb.String())
}

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

	return cfg, nil
}

// Config is the global collector configuration.
type Config struct {
	Cloud                 *Cloud `hcl:"cloud,block"`
	HTTPCollectorEndpoint string `hcl:"http_collector_endpoint,optional"`
	ConfigFile            string
}

// Cloud is the HCP Cloud configuration.
type Cloud struct {
	ClientID     string `hcl:"client_id,optional"`
	ClientSecret string `hcl:"client_secret,optional"`
	ResourceID   string `hcl:"resource_id,optional"`
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

	if c.IsEnabled() && (c.ClientID == "" || c.ClientSecret == "" || c.ResourceID == "") {
		return &errCloudConfigInvalidField{cloudCfg: *c}
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
