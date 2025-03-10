package oauth2

import (
	"context"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/dinhcanh303/go-microservices/pkg/constant"
	"golang.org/x/exp/slog"
	"golang.org/x/oauth2"
	facebookOAuth2 "golang.org/x/oauth2/facebook"
	googleOAuth2 "golang.org/x/oauth2/google"
)

// OAuthProviders is a struct that contains reference all the OAuth providers
type OAuthProvider struct {
	GoogleConfig   *oauth2.Config
	FacebookConfig *oauth2.Config
}

// OIDCProviders is a struct that contains reference all the OpenID providers
type OIDCProvider struct {
	GoogleOIDC *oidc.Provider
}

var (
	// OAuthProviders is a global variable that contains instance for all enabled the OAuth providers
	OAuthProviders OAuthProvider
	// OIDCProviders is a global variable that contains instance for all enabled the OpenID providers
	OIDCProviders OIDCProvider
)

// InitOAuth initializes the OAuth providers based on EnvData
func InitOAuth() error {
	ctx := context.Background()
	googleClientID, ok := os.LookupEnv(constant.KeyGoogleClientID)
	if !ok || len(googleClientID) == 0 {
		slog.Error("Environment variable not declared: %s", constant.KeyGoogleClientID)
	}
	googleClientSecret, ok := os.LookupEnv(constant.KeyGoogleClientSecret)
	if !ok || len(googleClientSecret) == 0 {
		slog.Error("Environment variable not declared: %s", constant.KeyGoogleClientSecret)
	}
	googleUrlCallback, ok := os.LookupEnv(constant.KeyGoogleUrlCallback)
	if !ok || len(googleClientSecret) == 0 {
		slog.Error("Environment variable not declared: %s", constant.KeyGoogleUrlCallback)
	}
	if googleClientID != "" && googleClientSecret != "" {
		p, err := oidc.NewProvider(ctx, "https://accounts.google.com")
		if err != nil {
			return err
		}
		OIDCProviders.GoogleOIDC = p
		OAuthProviders.GoogleConfig = &oauth2.Config{
			ClientID:     googleClientID,
			ClientSecret: googleClientSecret,
			RedirectURL:  googleUrlCallback,
			Endpoint:     googleOAuth2.Endpoint,
			Scopes:       []string{oidc.ScopeOpenID, "email", "profile"},
		}
	}
	facebookClientID, ok := os.LookupEnv(constant.KeyFacebookClientID)
	if !ok || len(facebookClientID) == 0 {
		slog.Error("Environment variable not declared: %s", constant.KeyFacebookClientID)
	}
	facebookClientSecret, ok := os.LookupEnv(constant.KeyFacebookClientSecret)
	if !ok || len(facebookClientSecret) == 0 {
		slog.Error("Environment variable not declared: %s", constant.KeyFacebookClientSecret)
	}
	if facebookClientID != "" && facebookClientSecret != "" {
		OAuthProviders.FacebookConfig = &oauth2.Config{
			ClientID:     facebookClientID,
			ClientSecret: facebookClientSecret,
			RedirectURL:  "/oauth_callback/facebook",
			Endpoint:     facebookOAuth2.Endpoint,
			Scopes:       []string{"public_profile", "email"},
		}
	}
	return nil
}
