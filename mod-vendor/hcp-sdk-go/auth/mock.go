// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package auth

import (
	"context"
	"time"

	"golang.org/x/oauth2"
)

type MockSession struct{}

// GetToken returns some mock token with static data.
func (s *MockSession) GetToken(ctx context.Context, conf *oauth2.Config) (*oauth2.Token, error) {
	tok := oauth2.Token{
		AccessToken:  "SomeAccessToken",
		RefreshToken: "SomeRefreshToken",
		Expiry:       time.Now().Add(time.Hour * 1),
	}
	return &tok, nil
}

// mockBrowser provides a mocked response of the OAuth2 login flow.
type mockBrowser struct{}

func (b *mockBrowser) GetTokenFromBrowser(ctx context.Context, conf *oauth2.Config) (*oauth2.Token, error) {
	tok := oauth2.Token{
		AccessToken:  "SomeNewAccessToken",
		RefreshToken: "SomeNewRefreshToken",
		Expiry:       time.Now().Add(time.Hour * 1),
	}
	return &tok, nil
}
