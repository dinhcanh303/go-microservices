package auth

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/dinhcanh303/go-microservices/internal/auth/domain"
	"github.com/dinhcanh303/go-microservices/internal/auth/usecases/keys"
	sharedkernel "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
	"github.com/dinhcanh303/go-microservices/pkg/ldap"
	"github.com/dinhcanh303/go-microservices/pkg/token"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
	"github.com/google/wire"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type service struct {
	repo       UserRepo
	ucKeys     keys.UseCase
	ldapClient ldap.LdapClient
	jwt        token.JWT
}

// SignIn implements UseCase.
func (s *service) SignIn(ctx context.Context, email string, password string) (*sharedkernel.UserAuth, error) {
	slog.Info("Service Auth:: SignIn")
	isEmailCompany := checkEmailCompany(email)
	if isEmailCompany {
		auth, _, err := s.ldapClient.Authenticate(email, password)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		if auth {
			foundUser, err := findUserByEmail(s.repo, ctx, email)
			if err != nil {
				return nil, err
			}
			return createTokenPairAndResponse(ctx, s, foundUser, false)
		}
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	foundUser, err := findUserByEmail(s.repo, ctx, email)
	if err != nil {
		return nil, err
	}
	match, err := utils.ComparePassword(foundUser.Password, password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	if !match {
		return nil, status.Error(codes.Unauthenticated, "Password mismatch")
	}
	return createTokenPairAndResponse(ctx, s, foundUser, false)
}

// SignUp implements UseCase.
func (s *service) SignUp(ctx context.Context, email, password, fistName, lastName string) (*sharedkernel.UserAuth, error) {
	slog.Info("Service Auth:: SignUp")
	isEmailCompany := checkEmailCompany(email)
	if isEmailCompany {
		return nil, status.Error(codes.AlreadyExists, fmt.Sprintf("Email has suffix is %s does't created. Please use the login function for this type of email!", suffixEmailCompany))
	}
	foundUser, _ := s.repo.GetUserByEmail(ctx, email)
	if foundUser != nil {
		return nil, status.Error(codes.AlreadyExists, "User already exists")
	}
	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}
	model := domain.NewUser(email, passwordHash, fistName, lastName)
	newUser, err := s.repo.CreateUser(ctx, model)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return createTokenPairAndResponse(ctx, s, newUser, true)
}

var _ UseCase = (*service)(nil)

func NewUseCase(
	repo UserRepo,
	ucKeys keys.UseCase,
	ldapClient ldap.LdapClient,
	jwt token.JWT,
) UseCase {
	return &service{
		repo:       repo,
		ucKeys:     ucKeys,
		ldapClient: ldapClient,
		jwt:        jwt,
	}
}

var UseCaseSet = wire.NewSet(NewUseCase)

func checkEmailCompany(email string) bool {
	return strings.Contains(email, suffixEmailCompany)
}
func findUserByEmail(repo UserRepo, ctx context.Context, email string) (*domain.User, error) {
	foundUser, err := repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, status.Error(codes.NotFound, "User not registered")
	}
	return foundUser, err
}
func createKeyPair() (string, string, error) {
	publicKey, err := utils.GenerateRandomHexBytes(64)
	if err != nil {
		return "", "", err
	}
	privateKey, err := utils.GenerateRandomHexBytes(64)
	if err != nil {
		return "", "", err
	}
	return publicKey, privateKey, nil
}

const (
	suffixEmailCompany = "@tlcmodular.com"
)

func createTokenPair(jwt token.JWT, user *domain.User, publicKey, privateKey string) (string, string, error) {
	payload := token.NewPayload(user.ID,
		user.Email,
		user.FullName,
		user.Role, "",
		2*24*time.Hour)
	accessToken, err := jwt.CreateToken(payload, publicKey)
	if err != nil {
		return "", "", err
	}
	payload = &token.Payload{
		ExpiredAt: payload.ExpiredAt.Add(5 * 24 * time.Hour),
	}
	refreshToken, err := jwt.CreateToken(payload, privateKey)
	if err != nil {
		return "", "", err
	}
	_, err = jwt.VerifyToken(accessToken, publicKey)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}
func createTokenPairAndResponse(ctx context.Context, service *service, user *domain.User, isSingUp bool) (*sharedkernel.UserAuth, error) {
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
	return &sharedkernel.UserAuth{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
