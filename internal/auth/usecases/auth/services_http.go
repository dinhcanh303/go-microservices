package auth

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/dinhcanh303/go-microservices/internal/auth/usecases/keys"
	"github.com/dinhcanh303/go-microservices/pkg/ldap"
	"github.com/dinhcanh303/go-microservices/pkg/oauth2"
	"github.com/dinhcanh303/go-microservices/pkg/redis"
	"github.com/dinhcanh303/go-microservices/pkg/token"
	"github.com/google/wire"
)

type serviceHttp struct {
	repo       UserRepo
	ucKeys     keys.UseCase
	ldapClient ldap.LdapClient
	jwt        token.JWT
	redis      redis.RedisEngine
}

var _ UseCaseHttp = (*serviceHttp)(nil)

func NewUseCaseHttp(
	repo UserRepo,
	ucKeys keys.UseCase,
	ldapClient ldap.LdapClient,
	jwt token.JWT,
	redis redis.RedisEngine,
) UseCaseHttp {
	return &serviceHttp{
		repo:       repo,
		ucKeys:     ucKeys,
		ldapClient: ldapClient,
		jwt:        jwt,
		redis:      redis,
	}
}

var UseCaseHttpSet = wire.NewSet(NewUseCaseHttp)
var (
	randomState = "state"
)

// GoogleCallback implements UseCaseHttp.
func (s *serviceHttp) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	slog.Info("Form", r.Form)
	if r.FormValue("stage") != randomState {
		slog.Error("failed state:...")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	token, err := oauth2.OAuthProviders.GoogleConfig.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		slog.Error("create token failed: %s", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		slog.Error("response failed: %s", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("read all content failed: %s", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	fmt.Fprintf(w, "Response: %s", content)
}

// GoogleLogin implements UseCaseHttp.
func (s *serviceHttp) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := oauth2.OAuthProviders.GoogleConfig.AuthCodeURL(randomState)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
