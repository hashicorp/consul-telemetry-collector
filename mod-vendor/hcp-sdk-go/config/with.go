// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"crypto/tls"
	"fmt"
	"net/url"

	"github.com/hashicorp/hcp-sdk-go/auth"
	"github.com/hashicorp/hcp-sdk-go/profile"
)

// WithClientCredentials credentials is an option that can be used to set
// HCP client credentials on the configuration.
func WithClientCredentials(clientID, clientSecret string) HCPConfigOption {
	return func(config *hcpConfig) error {
		config.clientCredentialsConfig.ClientID = clientID
		config.clientCredentialsConfig.ClientSecret = clientSecret

		return nil
	}
}

// WithAPI credentials is an option that can be used to provide a custom
// configuration for the API endpoint.
//
// If nil is provided for the tlsConfig value, TLS will be disabled.
//
// This should only be necessary for development purposes.
func WithAPI(address string, tlsConfig *tls.Config) HCPConfigOption {
	return func(config *hcpConfig) error {
		config.apiAddress = address
		config.apiTLSConfig = cloneTLSConfig(tlsConfig)

		return nil
	}
}

// WithSCADA credentials is an option that can be used to provide a custom
// configuration for the SCADA endpoint.
//
// If nil is provided for the tlsConfig value, TLS will be disabled.
//
// This should only be necessary for development purposes.
func WithSCADA(address string, tlsConfig *tls.Config) HCPConfigOption {
	return func(config *hcpConfig) error {
		config.scadaAddress = address
		config.scadaTLSConfig = cloneTLSConfig(tlsConfig)

		return nil
	}
}

// WithPortalURL credentials is an option that can be used to provide a custom
// URL for the portal.
//
// This should only be necessary for development purposes.
func WithPortalURL(portalURL string) HCPConfigOption {
	return func(config *hcpConfig) error {
		parsedPortalURL, err := url.Parse(portalURL)
		if err != nil {
			return fmt.Errorf("failed to parse portal URL: %w", err)
		}

		config.portalURL = parsedPortalURL

		return nil
	}
}

// WithAuth credentials is an option that can be used to provide a custom URL
// for the auth endpoint.
//
// An alternative TLS configuration can be provided, if none is provided the
// default TLS configuration will be used. It is not possible to disable TLS for
// the auth endpoint.
//
// This should only be necessary for development purposes.
func WithAuth(authURL string, tlsConfig *tls.Config) HCPConfigOption {
	return func(config *hcpConfig) error {
		parsedAuthURL, err := url.Parse(authURL)
		if err != nil {
			return fmt.Errorf("failed to parse auth URL: %w", err)
		}

		// Ensure a TLS configuration is set, as the auth endpoint should always
		// use TLS.
		if tlsConfig == nil {
			tlsConfig = &tls.Config{}
		}

		config.authURL = parsedAuthURL
		config.authTLSConfig = cloneTLSConfig(tlsConfig)

		// Ensure the OAuth2 endpoints are updated with the new auth URL
		config.oauth2Config.Endpoint.AuthURL = authURL + "/oauth2/auth"
		config.oauth2Config.Endpoint.TokenURL = authURL + "/oauth2/token"

		return nil
	}
}

// WithOAuth2ClientID credentials is an option that can be used to provide a custom OAuth2 Client ID.
//
// An alternative OAuth2 ClientID can be provided, if none is provided the
// default OAuth2 Client ID will be used.
//
// This should only be necessary for development purposes.
func WithOAuth2ClientID(oauth2ClientID string) HCPConfigOption {
	return func(config *hcpConfig) error {
		config.oauth2Config.ClientID = oauth2ClientID

		return nil
	}
}

// WithSession is an option that can be used to provide a custom Session struct.
//
// A mock Session can be provided, if none is provided the default UserSession
// will be used.
//
// This should only be necessary for testing purposes.
func WithSession(s auth.Session) HCPConfigOption {
	return func(config *hcpConfig) error {
		config.session = s

		return nil
	}
}

// WithProfile is an option that can be used to provide a custom UserProfile struct.
func WithProfile(p *profile.UserProfile) HCPConfigOption {
	return func(config *hcpConfig) error {
		config.profile = p
		return nil
	}
}
