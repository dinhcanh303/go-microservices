package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/dinhcanh303/go-microservices/internal/auth/domain"
	"github.com/dinhcanh303/go-microservices/internal/auth/usecases/keys"
	"github.com/dinhcanh303/go-microservices/pkg/constant"
	"github.com/dinhcanh303/go-microservices/pkg/ldap"
	"github.com/dinhcanh303/go-microservices/pkg/oauth2"
	"github.com/dinhcanh303/go-microservices/pkg/redis"
	"github.com/dinhcanh303/go-microservices/pkg/token"
	"github.com/google/wire"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	queryValues := r.URL.Query()
	state := queryValues.Get("state")
	if state != randomState {
		slog.Error("failed state:...")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	code := queryValues.Get("code")
	// os.Setenv("TLS_SKIP_VERIFY", "true")
	token, err := oauth2.OAuthProviders.GoogleConfig.Exchange(context.Background(), code)
	if err != nil {
		slog.Error("create token failed: ", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		slog.Error("response failed: ", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("read all content failed: ", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	user, err := registerUser(s, context.Background(), content)
	if err != nil {
		slog.Error("register user failed: ", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	fmt.Fprintf(w, "UserInfo: %s", user)
}

// GoogleLogin implements UseCaseHttp.
func (s *serviceHttp) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := oauth2.OAuthProviders.GoogleConfig.AuthCodeURL(randomState)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

type UserInfo struct {
	Email      string `json:"email"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Picture    string `json:"picture"`
}

func registerUser(s *serviceHttp, ctx context.Context, content []byte) ([]byte, error) {
	var userInfo UserInfo
	err := json.Unmarshal(content, &userInfo)
	if err != nil {
		slog.Error("Unmarshal content failed: ", err)
		return nil, err
	}
	foundUser, _ := s.repo.GetUserByEmail(ctx, userInfo.Email)
	if foundUser != nil {
		return createTokenPairAndResponse2(ctx, s, foundUser, false)
	}
	model := domain.NewUser(userInfo.Email, "password", userInfo.GivenName, userInfo.FamilyName, userInfo.Picture)
	newUser, err := s.repo.CreateUser(ctx, model)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	//Del cache
	err = s.redis.Invalidate(constant.CacheUsers)
	if err != nil {
		slog.Error("Invalidate cache list users failed")
	}
	return createTokenPairAndResponse2(ctx, s, newUser, true)
}
func createTokenPairAndResponse2(ctx context.Context, service *serviceHttp, user *domain.User, isSingUp bool) ([]byte, error) {
	publicKey, privateKey, err := createKeyPair()
	if err != nil {
		return nil, status.Error(codes.Unknown, "Create Public Private Key failed")
	}
	accessToken, refreshToken, err := createTokenPair(service.jwt, user, publicKey, privateKey)
	if err != nil {
		return nil, status.Error(codes.Unknown, "Create Token Failed")
	}
	model := &domain.Key{
		UserID:     user.ID,
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}
	if !isSingUp {
		model.RefreshToken = refreshToken
	}
	_, err = service.ucKeys.CreateKeyToken(ctx, model)
	if err != nil {
		return nil, status.Error(codes.Unknown, "Create Key Store Failed")
	}
	userJson, err := json.Marshal(domain.UserAuth{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
	if err != nil {
		return nil, status.Error(codes.Unknown, "Marshal UserAuth Failed")
	}
	return userJson, nil
}
