package oauth2

import (
	"context"
	"fmt"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/dinhcanh303/go-microservices/pkg/constant"
	"github.com/golang/glog"
	"golang.org/x/oauth2"
	facebookOAuth2 "golang.org/x/oauth2/facebook"
	githubOAuth2 "golang.org/x/oauth2/github"
	linkedInOAuth2 "golang.org/x/oauth2/linkedin"
	microsoftOAuth2 "golang.org/x/oauth2/microsoft"
	"google.golang.org/appengine/log"
)

const (
	microsoftCommonTenant = "common"
)

// OAuthProviders is a struct that contains reference all the OAuth providers
type OAuthProvider struct {
	GoogleConfig    *oauth2.Config
	GithubConfig    *oauth2.Config
	FacebookConfig  *oauth2.Config
	LinkedInConfig  *oauth2.Config
	AppleConfig     *oauth2.Config
	TwitterConfig   *oauth2.Config
	MicrosoftConfig *oauth2.Config
}

// OIDCProviders is a struct that contains reference all the OpenID providers
type OIDCProvider struct {
	GoogleOIDC    *oidc.Provider
	MicrosoftOIDC *oidc.Provider
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
		glog.Fatalf("environment variable not declared: GOOGLE_CLIENT_ID")
	}
	googleClientSecret, ok := os.LookupEnv(constant.KeyGoogleClientSecret)
	if !ok || len(googleClientSecret) == 0 {
		glog.Fatalf("environment variable not declared: GOOGLE_CLIENT_SECRET")
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
			RedirectURL:  "/oauth_callback/google",
			Endpoint:     OIDCProviders.GoogleOIDC.Endpoint(),
			Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
		}
	}
	githubClientID, ok := os.LookupEnv(constant.KeyGithubClientID)
	if !ok || len(googleClientID) == 0 {
		glog.Fatalf("environment variable not declared: GOOGLE_CLIENT_ID")
	}
	githubClientSecret, ok := os.LookupEnv(constant.KeyGithubClientSecret)
	if !ok || len(googleClientSecret) == 0 {
		glog.Fatalf("environment variable not declared: GOOGLE_CLIENT_SECRET")
	}
	if githubClientID != "" && githubClientSecret != "" {
		OAuthProviders.GithubConfig = &oauth2.Config{
			ClientID:     githubClientID,
			ClientSecret: githubClientSecret,
			RedirectURL:  "/oauth_callback/github",
			Endpoint:     githubOAuth2.Endpoint,
			Scopes:       []string{"read:user", "user:email"},
		}
	}

	facebookClientID, ok := os.LookupEnv(constant.KeyFacebookClientID)
	if !ok || len(facebookClientID) == 0 {
		glog.Fatalf("environment variable not declared: GOOGLE_CLIENT_ID")
	}
	facebookClientSecret, ok := os.LookupEnv(constant.KeyFacebookClientSecret)
	if !ok || len(facebookClientSecret) == 0 {
		glog.Fatalf("environment variable not declared: GOOGLE_CLIENT_SECRET")
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

	linkedInClientID, ok := os.LookupEnv(constant.KeyLinkedInClientID)
	if !ok || len(linkedInClientID) == 0 {
		glog.Fatalf("environment variable not declared: GOOGLE_CLIENT_ID")
	}
	linkedInClientSecret, ok := os.LookupEnv(constant.KeyLinkedInClientSecret)
	if !ok || len(linkedInClientSecret) == 0 {
		glog.Fatalf("environment variable not declared: GOOGLE_CLIENT_SECRET")
	}
	if linkedInClientID != "" && linkedInClientSecret != "" {
		OAuthProviders.LinkedInConfig = &oauth2.Config{
			ClientID:     linkedInClientID,
			ClientSecret: linkedInClientSecret,
			RedirectURL:  "/oauth_callback/linkedin",
			Endpoint:     linkedInOAuth2.Endpoint,
			Scopes:       []string{"r_liteprofile", "r_emailaddress"},
		}
	}
	appleClientID, ok := os.LookupEnv(constant.KeyAppleClientID)
	if !ok || len(appleClientID) == 0 {
		glog.Fatalf("environment variable not declared: GOOGLE_CLIENT_ID")
	}
	appleClientSecret, ok := os.LookupEnv(constant.KeyAppleClientSecret)
	if !ok || len(appleClientSecret) == 0 {
		glog.Fatalf("environment variable not declared: GOOGLE_CLIENT_SECRET")
	}
	if appleClientID != "" && appleClientSecret != "" {
		OAuthProviders.AppleConfig = &oauth2.Config{
			ClientID:     appleClientID,
			ClientSecret: appleClientSecret,
			RedirectURL:  "/oauth_callback/apple",
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://appleid.apple.com/auth/authorize",
				TokenURL: "https://appleid.apple.com/auth/token",
			},
		}
	}
	twitterClientID, ok := os.LookupEnv(constant.KeyTwitterClientID)
	if !ok || len(twitterClientID) == 0 {
		glog.Fatalf("environment variable not declared: GOOGLE_CLIENT_ID")
	}
	twitterClientSecret, ok := os.LookupEnv(constant.KeyTwitterClientSecret)
	if !ok || len(twitterClientSecret) == 0 {
		glog.Fatalf("environment variable not declared: GOOGLE_CLIENT_SECRET")
	}
	if twitterClientID != "" && twitterClientSecret != "" {
		OAuthProviders.TwitterConfig = &oauth2.Config{
			ClientID:     twitterClientID,
			ClientSecret: twitterClientSecret,
			RedirectURL:  "/oauth_callback/twitter",
			Endpoint: oauth2.Endpoint{
				// Endpoint is currently not yet part of oauth2-package. See https://go-review.googlesource.com/c/oauth2/+/350889 for status
				AuthURL:   "https://twitter.com/i/oauth2/authorize",
				TokenURL:  "https://api.twitter.com/2/oauth2/token",
				AuthStyle: oauth2.AuthStyleInHeader,
			},
			Scopes: []string{"tweet.read", "users.read"},
		}
	}
	microsoftClientID, ok := os.LookupEnv(constant.KeyMicrosoftClientID)
	if !ok || len(microsoftClientID) == 0 {
		glog.Fatalf("environment variable not declared: GOOGLE_CLIENT_ID")
	}
	microsoftClientSecret, ok := os.LookupEnv(constant.KeyMicrosoftClientSecret)
	if !ok || len(microsoftClientSecret) == 0 {
		glog.Fatalf("environment variable not declared: GOOGLE_CLIENT_SECRET")
	}
	microsoftActiveDirTenantID, ok := os.LookupEnv(constant.KeyMicrosoftActiveDirectoryTenantID)
	if !ok || len(microsoftActiveDirTenantID) == 0 {
		glog.Fatalf("environment variable not declared: GOOGLE_CLIENT_SECRET")
		microsoftActiveDirTenantID = microsoftCommonTenant
	}
	if microsoftClientID != "" && microsoftClientSecret != "" {
		if microsoftActiveDirTenantID == microsoftCommonTenant {
			ctx = oidc.InsecureIssuerURLContext(ctx, fmt.Sprintf("https://login.microsoftonline.com/%s/v2.0", microsoftActiveDirTenantID))
		}
		p, err := oidc.NewProvider(ctx, fmt.Sprintf("https://login.microsoftonline.com/%s/v2.0", microsoftActiveDirTenantID))
		if err != nil {
			log.Debugf(ctx, "Error while creating OIDC provider for Microsoft: %v", err)
			return err
		}
		OIDCProviders.MicrosoftOIDC = p
		OAuthProviders.MicrosoftConfig = &oauth2.Config{
			ClientID:     microsoftClientID,
			ClientSecret: microsoftClientSecret,
			RedirectURL:  "/oauth_callback/microsoft",
			Endpoint:     microsoftOAuth2.AzureADEndpoint(microsoftActiveDirTenantID),
			Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
		}
	}

	return nil
}
