package controllers

import (
	"context"
	"github.com/coreos/go-oidc"
	"github.com/revel/revel"
	"golang.org/x/oauth2"
	"strings"
)

type Auth struct {
	*revel.Controller
}
func (c Auth) Callback() revel.Result {
	configURL := "https://auth.chalamov.dev/auth/realms/test"
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, configURL)
	if err != nil {
		panic(err)
	}

	clientID := "test-client"
	clientSecret := "540fb896-a614-4932-b802-327b8696e663"

	redirectURL := "http://localhost:5010/login/callback"
	// Configure an OpenID Connect aware OAuth2 client.
	oauth2Config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),
		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
	}

	if c.Request.URL.Query().Get("state") != "somestate" {
		return c.NotFound("40011")
	}

	oauth2Token, err := oauth2Config.Exchange(ctx, c.Request.URL.Query().Get("code"))
	if err != nil {
		return c.NotFound( "Failed to exchange token: "+err.Error())
	}

	return c.RenderJSON(oauth2Token.AccessToken)
}

func (c Auth) Login() revel.Result {
	configURL := "https://auth.chalamov.dev/auth/realms/test"
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, configURL)
	if err != nil {
		panic(err)
	}

	clientID := "test-client"
	clientSecret := "540fb896-a614-4932-b802-327b8696e663"

	redirectURL := "http://localhost:5010/login	/callback"
	// Configure an OpenID Connect aware OAuth2 client.
	oauth2Config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),
		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
	}

	oidcConfig := &oidc.Config{
		ClientID: clientID,
	}
	verifier := provider.Verifier(oidcConfig)

	rawAccessToken := c.Request.Header.Get("Authorization");
	if rawAccessToken == "" {
		return c.Redirect(oauth2Config.AuthCodeURL("somestate"))
	}
	parts := strings.Split(rawAccessToken, " ")
	if len(parts) != 2 {
		return c.NotFound("400")
	}
	_, err2 := verifier.Verify(ctx, parts[1])

	if err2 != nil {
		return c.NotFound("402")
	}

	return c.Render()
}
