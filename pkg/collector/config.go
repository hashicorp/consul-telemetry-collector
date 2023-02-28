package collector

type Config struct {
	Cloud                 Cloud  `hcl:"cloud"`
	HTTPCollectorEndpoint string `hcl:"http_collector_endpoint"`
}

type Cloud struct {
	ClientID     string `hcl:"client_id"`
	ClientSecret string `hcl:"client_secret"`
	ResourceID   string `hcl:"resource_id"`
}

func (c Cloud) IsEnabled() bool {
	if c.ClientSecret == "" && c.ClientID == "" && c.ResourceID == "" {
		return false
	}
	return true
}
