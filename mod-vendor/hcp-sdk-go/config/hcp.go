// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"crypto/tls"
	"fmt"
	"net/url"

	"github.com/hashicorp/hcp-sdk-go/auth"
	"github.com/hashicorp/hcp-sdk-go/profile"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// HCPConfig provides configuration values that are useful to interact with HCP.
type HCPConfig interface {

	// Profile will return the user's configured profile information
	Profile() *profile.UserProfile

	// TokenSource will return a token that can be used to authenticate to HCP.
	oauth2.TokenSource

	// PortalURL will return the URL of the portal.
	//
	// The default portal URL will be the production portal, but the value may
	// be overwritten for development purposes.
	PortalURL() *url.URL

	// APIAddress will return the HCP API address (<hostname>[:port]).
	//
	// The default address will be the HCP production API, but it may be
	// overwritten for development purposes.
	APIAddress() string

	// APITLSConfig will return the API TLS configuration.
	//
	// TLS will be enabled by default but may be disabled for development
	// purposes.
	APITLSConfig() *tls.Config

	// SCADAAddress will return the SCADA address (<hostname>[:port]).
	//
	// The default address will be the HCP production SCADA endpoint, but it may
	// be overwritten for development purposes.
	SCADAAddress() string

	// SCADATLSConfig will return the SCADA TLS configuration.
	//
	// TLS will be enabled by default but may be disabled for development
	// purposes.
	SCADATLSConfig() *tls.Config
}

type hcpConfig struct {
	// clientCredentialsConfig is the configuration that will be used to create
	// the token source.
	clientCredentialsConfig clientcredentials.Config

	// oauth2Config is the configuration that will be used to create
	// a browser-initiated token source when client credentials are not provided.
	oauth2Config oauth2.Config

	// authURL is the URL that will be used to authenticate.
	authURL *url.URL
	// authTLSConfig is the TLS configuration for the auth endpoint. TLS can not
	// be disabled for the auth endpoint, but the configuration can be set to a
	// custom one for tests or development.
	authTLSConfig *tls.Config

	// tokenSource is a self-refreshing token source that returns an access
	// token that can be used to authenticate against the HCP APIs.
	tokenSource oauth2.TokenSource

	// portalURL is the base URL of the portal.
	portalURL *url.URL

	// apiAddress is the address (<hostname>[:port]) of the HCP API.
	apiAddress string
	// apiTLSConfig is the TLS configuration for the HCP API. It can be set to
	// nil if TLS should be disabled.
	apiTLSConfig *tls.Config

	// scadaAddress is the address (<hostname>[:port]) of the SCADA endpoint.
	scadaAddress string
	// scadaTLSConfig is the TLS configuration for the SCADA endpoint. It can be
	// set to nil if TLS should be disabled.
	scadaTLSConfig *tls.Config

	// session is responsible for getting an access token fron our identity provider.
	// A mock can be used in tests.
	session auth.Session

	// profile is the user's organization id and project id
	profile *profile.UserProfile
}

func (c *hcpConfig) Profile() *profile.UserProfile {
	return c.profile
}

func (c *hcpConfig) Token() (*oauth2.Token, error) {
	return c.tokenSource.Token()
}

func (c *hcpConfig) PortalURL() *url.URL {
	// Copy the value in order to not return a pointer to the internal one.
	portalURL := *c.portalURL

	return &portalURL
}

func (c *hcpConfig) APIAddress() string {
	return c.apiAddress
}

func (c *hcpConfig) APITLSConfig() *tls.Config {
	return cloneTLSConfig(c.apiTLSConfig)
}

func (c *hcpConfig) SCADAAddress() string {
	return c.scadaAddress
}

func (c *hcpConfig) SCADATLSConfig() *tls.Config {
	return cloneTLSConfig(c.scadaTLSConfig)
}

func (c *hcpConfig) validate() error {

	// Ensure both client credentials provided
	if (c.clientCredentialsConfig.ClientID == "" && c.clientCredentialsConfig.ClientSecret != "") ||
		(c.clientCredentialsConfig.ClientID != "" && c.clientCredentialsConfig.ClientSecret == "") {
		return fmt.Errorf("both client ID and secret must be provided")
	}

	// Ensure at least one auth method configured
	if c.clientCredentialsConfig.ClientID == "" && c.clientCredentialsConfig.ClientSecret == "" && c.oauth2Config.ClientID == "" {
		return fmt.Errorf("either client credentials or oauth2 client ID must be provided")
	}

	// Ensure profile contains both org ID and project ID
	if (c.profile.OrganizationID == "" && c.profile.ProjectID != "") ||
		(c.profile.OrganizationID != "" && c.profile.ProjectID == "") {
		return fmt.Errorf("when setting a user profile, both organization ID and project ID must be provided")
	}

	// Ensure the auth URL is valid
	if c.authURL.Host == "" {
		return fmt.Errorf("the auth URL has to be non-empty")
	}

	if c.authURL.Scheme != "https" {
		return fmt.Errorf("the auth URL has to use https at its scheme")
	}

	// Ensure the portal URL is valid
	if c.portalURL.Host == "" {
		return fmt.Errorf("the portal URL has to be non-empty")
	}

	// Ensure an API address is valid
	if c.apiAddress == "" {
		return fmt.Errorf("the API address has to be non-empty")
	}

	// Ensure an SCADA address is valid
	if c.scadaAddress == "" {
		return fmt.Errorf("the SCADA address has to be non-empty")
	}

	return nil
}

var _ HCPConfig = &hcpConfig{}
