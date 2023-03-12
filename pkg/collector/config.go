package collector

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

var (
	errNoConfigurationProvided = errors.New("no configuration provided: see usage")
	errNoCollectorEndpoint     = fmt.Errorf("collector endpoint must be set with flag: %s or env-var: %s", COOtelHTTPEndpointOpt, COOtelHTTPEndpoint)
	errCloudConfigInvalid      = fmt.Errorf("cloud configuration is not valid")
)

// used to parse file path and return a configuration
type parser func(*string) (*Config, error)

// ParseFile parses the given file for a configuration. The file is expected
// to be in JSON format.
func parseFile(filename *string) (*Config, error) {
	f, err := os.Open(*filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cfg := &Config{}

	// wrap our file in a 1mb limit reader
	mb := int64(1024 * 1024 * 1024)
	r := io.LimitReader(f, mb)
	buffer := bytes.NewBuffer(make([]byte, 0, mb))
	_, err = io.Copy(buffer, r)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", *filename, err)
	}
	// decode needs filename for parsing and bytes passed to it
	// NOTE we should probably support both
	err = hclsimple.Decode(*filename, buffer.Bytes(), nil, cfg)

	if err != nil {
		return nil, fmt.Errorf("failed parsing hcl config file: %w", err)
	}

	return cfg, nil
}

// Config is the global collector configuration
type Config struct {
	Cloud                 *Cloud  `hcl:"cloud,block"`
	HTTPCollectorEndpoint *string `hcl:"http_collector_endpoint,optional"`
	ConfigFile            *string
}

// Cloud is the HCP Cloud configuration
type Cloud struct {
	ClientID     *string `hcl:"client_id,optional"`
	ClientSecret *string `hcl:"client_secret,optional"`
	ResourceID   *string `hcl:"resource_id,optional"`
}

// IsEnabled checks if the Cloud config is enabled. It returns false if the ClientID,
// ClientSecret and ResourceID are all empty
func (c *Cloud) IsEnabled() bool {
	if c == nil {
		return false
	}
	if *c.ClientSecret == "" && *c.ClientID == "" && *c.ResourceID == "" {
		return false
	}
	return true
}

func (c *Cloud) validate() error {
	if c == nil {
		return nil
	}
	if (!empty(c.ClientID) && (empty(c.ClientSecret) || empty(c.ResourceID))) ||
		(!empty(c.ClientSecret) && (empty(c.ClientID) || empty(c.ResourceID))) ||
		(!empty(c.ResourceID) && (empty(c.ClientID) || empty(c.ClientSecret))) {
		return errCloudConfigInvalid
	}
	return nil
}

func empty(s *string) bool {
	if s == nil {
		return true
	}
	if *s == "" {
		return true
	}
	return false
}

func (c *Config) validate() error {
	if c == nil {
		return errNoConfigurationProvided
	}

	if c.HTTPCollectorEndpoint == nil {
		return errNoCollectorEndpoint
	}

	if c.Cloud == nil {
		return nil
	}

	return c.Cloud.validate()
}
