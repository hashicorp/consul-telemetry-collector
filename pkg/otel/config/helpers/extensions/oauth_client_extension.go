package extensions

import (
	"fmt"
	"net/url"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configopaque"

	oauth "github.com/open-telemetry/opentelemetry-collector-contrib/extension/oauth2clientauthextension"
)

const (
	defaultIssuerURL = "https://auth.idp.hashicorp.com"
	audienceKey      = "audience"
	audience         = "https://api.hashicorp.com"
	oauth2ClientName = "oauth2client"
)

// Oauth2ClientID is the component.ID used by the oauth2client extension
var Oauth2ClientID component.ID = component.NewIDWithName(oauth2ClientName, "hcp")

// OauthClientCfg returns a component ID and oauth config
func OauthClientCfg(clientID string, clientSecret string) (component.ID, *oauth.Config) {
	return Oauth2ClientID, &oauth.Config{
		ClientID:       clientID,
		ClientSecret:   configopaque.String(clientSecret),
		TokenURL:       fmt.Sprintf("%s/oauth2/token", defaultIssuerURL),
		EndpointParams: url.Values{audienceKey: []string{audience}},
	}
}
