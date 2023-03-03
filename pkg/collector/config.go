package collector

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
	ResourceID   string `hcl:"resource_id,optional"`
}

// IsEnabled checks if the Cloud config is enabled. It returns false if the ClientID,
// ClientSecret and ResourceID are all empty
func (c *Cloud) IsEnabled() bool {
	if c == nil {
		return false
	}
	if c.ClientSecret == "" && c.ClientID == "" && c.ResourceID == "" {
		return false
	}
	return true
}

func (c *Cloud) clone() *Cloud {
	if c == nil {
		return &Cloud{}
	}
	return &Cloud{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		ResourceID:   c.ResourceID,
	}
}

func (c Config) clone() Config {
	return Config{
		Cloud:                 c.Cloud.clone(),
		HTTPCollectorEndpoint: c.HTTPCollectorEndpoint,
		ConfigFile:            c.ConfigFile,
	}
}
