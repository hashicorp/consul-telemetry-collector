package confhelper

import (
	"fmt"
	"os"

	"go.opentelemetry.io/collector/component"

	"github.com/hashicorp/consul-telemetry-collector/pkg/confresolver"
)

const defaultIssuerURL = "https://auth.idp.hashicorp.com"
const audience = "https://api.hashicorp.com"

// Oauth2ClientID is the component.ID used by the oauth2client extension
const Oauth2ClientID = "oauth2client"

// OauthClient helper creates an oauth2client authentication extension that authenticates to HCP.
func OauthClient(c *confresolver.Config, clientID, clientSecret string) {
	// this duplicates logic in hcp-sdk-go
	var issuerURL string
	var ok bool
	if issuerURL, ok = os.LookupEnv("HCP_AUTH_URL"); !ok {
		issuerURL = defaultIssuerURL
	}

	const (
		clientIDKey     = "client_id"
		clientSecretKey = "client_secret"
		tokenURL        = "token_url"
		audienceKey     = "audience"
	)

	oauth2auth := c.NewExtensions(component.NewIDWithName(Oauth2ClientID, "hcp"))
	oauth2auth.Set(clientIDKey, clientID)
	oauth2auth.Set(clientSecretKey, clientSecret)
	oauth2auth.Set(tokenURL, fmt.Sprintf("%s/oauth2/token", issuerURL))
	endpointParams := oauth2auth.SetMap("endpoint_params")
	endpointParams.Set(audienceKey, audience)
}
