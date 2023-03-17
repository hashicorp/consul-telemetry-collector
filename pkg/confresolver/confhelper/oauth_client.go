package confhelper

import (
	"fmt"
	"os"

	"go.opentelemetry.io/collector/component"

	"github.com/hashicorp/consul-telemetry-collector/pkg/confresolver"
)

const defaultIssuerURL = "https://auth.idp.hashicorp.com"
const audience = "https://api.hashicorp.com"

func OauthClient(c *confresolver.Config, clientID, clientSecret string) {
	// this duplicates logic in hcp-sdk-go
	var issuerURL string
	var ok bool
	if issuerURL, ok = os.LookupEnv("HCP_AUTH_URL"); !ok {
		issuerURL = defaultIssuerURL
	}

	oauth2auth := c.NewExtensions(component.NewIDWithName("oauth2client", "hcp"))
	oauth2auth.Set("client_id", clientID)
	oauth2auth.Set("client_secret", clientSecret)
	oauth2auth.Set("token_url", fmt.Sprintf("%s/oauth2/token", issuerURL))
	endpointParams := oauth2auth.SetMap("endpoint_params")
	endpointParams.Set("audience", audience)
}
