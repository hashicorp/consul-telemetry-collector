package collector

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

var (
	errNoConfigurationProvided = errors.New("no configuration provided: see usage")
	errCloudConfigInvalid      = fmt.Errorf("cloud configuration is not valid")
)

func configFromEnvVars() *Config {
	return &Config{
		Cloud: &Cloud{
			ClientID:     os.Getenv(HCPClientID),
			ClientSecret: os.Getenv(HCPClientSecret),
			ResourceURL:  os.Getenv(HCPResourceURL),
		},
		ConfigFile:            os.Getenv(COOConfigPath),
		HTTPCollectorEndpoint: os.Getenv(COOtelHTTPEndpoint),
	}

}

// used to parse a file path and return a configuration
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
	if err != nil {
		f.Close()
		return nil, err
	}
	f.Close()
	return cfg, nil
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

// Config is the global collector configuration
type Config struct {
	Cloud                 *Cloud `hcl:"cloud,block"`
	HTTPCollectorEndpoint string `hcl:"http_collector_endpoint,optional"`
	ConfigFile            string
}

// Cloud is the HCP Cloud configuration
type Cloud struct {
	ClientID     string `hcl:"client_id,optional"`
	ClientSecret string `hcl:"client_secret,optional"`
	ResourceURL  string `hcl:"resource_url,optional"`
}

// IsEnabled checks if the Cloud config is enabled. It returns false if the ClientID,
// ClientSecret and ResourceURL are all empty
func (c *Cloud) IsEnabled() bool {
	if c == nil {
		return false
	}
	if c.ClientSecret != "" || c.ClientID != "" || c.ResourceURL != "" {
		return true
	}
	return false
}

func (c *Cloud) validate() error {
	if c == nil {
		return nil
	}
	if c.IsEnabled() && (c.ClientID == "" || c.ClientSecret == "" || c.ResourceURL == "") {
		return errCloudConfigInvalid
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
