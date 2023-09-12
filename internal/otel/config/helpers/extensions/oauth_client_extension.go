// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package extensions

import (
	"fmt"
	"net/url"
	"os"

	"go.opentelemetry.io/collector/component"

	"github.com/hashicorp/consul-telemetry-collector/internal/otel/config/helpers/types"
)

const (
	defaultIssuerURL = "https://auth.idp.hashicorp.com"
	audienceKey      = "audience"
	defaultAudience  = "https://api.hashicorp.cloud"
	oauth2ClientName = "oauth2client"

	envVarAuthURL = "HCP_AUTH_URL"
	envVarAuthTLS = "HCP_AUTH_TLS"

	tlsSettingInsecure = "insecure"
	tlsSettingDisabled = "disabled"
)

// OauthClientID is the component.ID used by the oauth2client extension.
var OauthClientID component.ID = component.NewIDWithName(oauth2ClientName, "hcp")

// OauthClientConfig is a base wrapper around the oauth2clientauthextension.Config which
// we cannot use directly since the opaque client secret string gets changed to REDACTED when unmarshalling
//
//	https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/extension/oauth2clientauthextension
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

	// TLSSetting struct exposes TLS client configuration for the underneath client to authorization server.
	TLSSetting types.TLSClientSetting `mapstructure:"tls,omitempty"`
}

// OauthClientCfg returns a component ID and oauth config.
func OauthClientCfg(clientID, clientSecret string) *OauthClientConfig {
	authURL, ok := os.LookupEnv(envVarAuthURL)
	if !ok {
		authURL = defaultIssuerURL
	}

	authTLSConfig := tlsConfigForSetting()

	return &OauthClientConfig{
		ClientID:       clientID,
		ClientSecret:   clientSecret,
		TokenURL:       fmt.Sprintf("%s/oauth2/token", authURL),
		EndpointParams: url.Values{audienceKey: []string{defaultAudience}},
		TLSSetting:     authTLSConfig,
	}
}

func tlsConfigForSetting() types.TLSClientSetting {
	setting := os.Getenv(envVarAuthTLS)
	switch setting {
	case tlsSettingDisabled:
		return types.TLSClientSetting{Insecure: true}
	case tlsSettingInsecure:
		return types.TLSClientSetting{InsecureSkipVerify: true}
	default:
		return types.TLSClientSetting{}
	}
}
