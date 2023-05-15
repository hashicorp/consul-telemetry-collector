package extensions

import (
	"fmt"
	"net/url"

	"go.opentelemetry.io/collector/component"
)

const (
	defaultIssuerURL = "https://auth.idp.hashicorp.com"
	audienceKey      = "audience"
	audience         = "https://api.hashicorp.com"
	oauth2ClientName = "oauth2client"
)

// OauthClientID is the component.ID used by the oauth2client extension.
var OauthClientID component.ID = component.NewIDWithName(oauth2ClientName, "hcp")

// OauthClientConfig is a base wrapper around the oauth2clientauthextension.Config which
// we cannot use directly since the opaque client secret string gets changed to REDACTED when unmarshalling
//
//	github.com/open-telemetry/opentelemetry-collector-contrib/extension/oauth2clientauthextension
type OauthClientConfig struct {

	// ClientID is the application's ID.
	ClientID string `mapstructure:"client_id"`

	// ClientSecret is the application's secret.
	ClientSecret string `mapstructure:"client_secret"`

	// EndpointParams specifies additional parameters for requests to the token endpoint.
	EndpointParams url.Values `mapstructure:"endpoint_params"`

	// TokenURL is the resource server's token endpoint
	// URL. This is a constant specific to each server.
	// See https://datatracker.ietf.org/doc/html/rfc6749#section-3.2
	TokenURL string `mapstructure:"token_url"`
}

// OauthClientCfg returns a component ID and oauth config.
func OauthClientCfg(clientID, clientSecret string) *OauthClientConfig {
	return &OauthClientConfig{
		ClientID:       clientID,
		ClientSecret:   clientSecret,
		TokenURL:       fmt.Sprintf("%s/oauth2/token", defaultIssuerURL),
		EndpointParams: url.Values{audienceKey: []string{audience}},
	}
}
