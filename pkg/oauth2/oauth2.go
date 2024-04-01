package oauth2

import (
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/dinhcanh303/go-microservices/pkg/constant"
	"github.com/golang/glog"
	"golang.org/x/oauth2"
	githubOAuth2 "golang.org/x/oauth2/github"
	googleOAuth2 "golang.org/x/oauth2/google"
)

// OAuthProviders is a struct that contains reference all the OAuth providers
type OAuthProvider struct {
	GoogleConfig   *oauth2.Config
	GithubConfig   *oauth2.Config
	FacebookConfig *oauth2.Config
	LinkedInConfig *oauth2.Config
	AppleConfig    *oauth2.Config
	TwitterConfig  *oauth2.Config
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
	// ctx := context.Background()
	googleClientID, ok := os.LookupEnv(constant.KeyGoogleClientID)
	if !ok || len(googleClientID) == 0 {
		glog.Fatalf("environment variable not declared: %s", constant.KeyGoogleClientID)
	}
	googleClientSecret, ok := os.LookupEnv(constant.KeyGoogleClientSecret)
	if !ok || len(googleClientSecret) == 0 {
		glog.Fatalf("environment variable not declared: %s", constant.KeyGithubClientSecret)
	}
	googleUrlCallback, ok := os.LookupEnv(constant.KeyGoogleUrlCallback)
	if !ok || len(googleClientSecret) == 0 {
		glog.Fatalf("environment variable not declared: %s", constant.KeyGoogleUrlCallback)
	}
	if googleClientID != "" && googleClientSecret != "" {
		// p, err := oidc.NewProvider(ctx, "https://accounts.google.com")
		// if err != nil {
		// 	return err
		// }
		// OIDCProviders.GoogleOIDC = p
		// OAuthProviders.GoogleConfig = &oauth2.Config{
		// 	ClientID:     googleClientID,
		// 	ClientSecret: googleClientSecret,
		// 	RedirectURL:  "/oauth_callback/google",
		// 	Endpoint:     OIDCProviders.GoogleOIDC.Endpoint(),
		// 	Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
		// }

		OAuthProviders.GoogleConfig = &oauth2.Config{
			ClientID:     googleClientID,
			ClientSecret: googleClientSecret,
			RedirectURL:  googleUrlCallback,
			Endpoint:     googleOAuth2.Endpoint,
			Scopes:       []string{"email", "profile"},
		}
	}
	githubClientID, ok := os.LookupEnv(constant.KeyGithubClientID)
	if !ok || len(googleClientID) == 0 {
		glog.Fatalf("environment variable not declared: %s", constant.KeyGithubClientID)
	}
	githubClientSecret, ok := os.LookupEnv(constant.KeyGithubClientSecret)
	if !ok || len(googleClientSecret) == 0 {
		glog.Fatalf("environment variable not declared: %s", constant.KeyGithubClientSecret)
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

	// facebookClientID, ok := os.LookupEnv(constant.KeyFacebookClientID)
	// if !ok || len(facebookClientID) == 0 {
	// 	glog.Fatalf("environment variable not declared: %s", constant.KeyFacebookClientID)
	// }
	// facebookClientSecret, ok := os.LookupEnv(constant.KeyFacebookClientSecret)
	// if !ok || len(facebookClientSecret) == 0 {
	// 	glog.Fatalf("environment variable not declared: %s", constant.KeyFacebookClientSecret)
	// }
	// if facebookClientID != "" && facebookClientSecret != "" {
	// 	OAuthProviders.FacebookConfig = &oauth2.Config{
	// 		ClientID:     facebookClientID,
	// 		ClientSecret: facebookClientSecret,
	// 		RedirectURL:  "/oauth_callback/facebook",
	// 		Endpoint:     facebookOAuth2.Endpoint,
	// 		Scopes:       []string{"public_profile", "email"},
	// 	}
	// }

	// linkedInClientID, ok := os.LookupEnv(constant.KeyLinkedInClientID)
	// if !ok || len(linkedInClientID) == 0 {
	// 	glog.Fatalf("environment variable not declared: %s", constant.KeyLinkedInClientID)
	// }
	// linkedInClientSecret, ok := os.LookupEnv(constant.KeyLinkedInClientSecret)
	// if !ok || len(linkedInClientSecret) == 0 {
	// 	glog.Fatalf("environment variable not declared: %s", constant.KeyLinkedInClientSecret)
	// }
	// if linkedInClientID != "" && linkedInClientSecret != "" {
	// 	OAuthProviders.LinkedInConfig = &oauth2.Config{
	// 		ClientID:     linkedInClientID,
	// 		ClientSecret: linkedInClientSecret,
	// 		RedirectURL:  "/oauth_callback/linkedin",
	// 		Endpoint:     linkedInOAuth2.Endpoint,
	// 		Scopes:       []string{"r_liteprofile", "r_emailaddress"},
	// 	}
	// }
	// twitterClientID, ok := os.LookupEnv(constant.KeyTwitterClientID)
	// if !ok || len(twitterClientID) == 0 {
	// 	glog.Fatalf("environment variable not declared: %s", constant.KeyTwitterClientID)
	// }
	// twitterClientSecret, ok := os.LookupEnv(constant.KeyTwitterClientSecret)
	// if !ok || len(twitterClientSecret) == 0 {
	// 	glog.Fatalf("environment variable not declared: %s", constant.KeyTwitterClientSecret)
	// }
	// if twitterClientID != "" && twitterClientSecret != "" {
	// 	OAuthProviders.TwitterConfig = &oauth2.Config{
	// 		ClientID:     twitterClientID,
	// 		ClientSecret: twitterClientSecret,
	// 		RedirectURL:  "/oauth_callback/twitter",
	// 		Endpoint: oauth2.Endpoint{
	// 			// Endpoint is currently not yet part of oauth2-package. See https://go-review.googlesource.com/c/oauth2/+/350889 for status
	// 			AuthURL:   "https://twitter.com/i/oauth2/authorize",
	// 			TokenURL:  "https://api.twitter.com/2/oauth2/token",
	// 			AuthStyle: oauth2.AuthStyleInHeader,
	// 		},
	// 		Scopes: []string{"tweet.read", "users.read"},
	// 	}
	// }
	return nil
}
