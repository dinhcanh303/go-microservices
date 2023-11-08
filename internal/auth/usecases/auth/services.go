package auth

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/dinhcanh303/go-microservices/internal/auth/domain"
	"github.com/dinhcanh303/go-microservices/internal/auth/usecases/keys"
	sharedkernel "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
	"github.com/dinhcanh303/go-microservices/pkg/ldap"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
	"github.com/google/wire"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type service struct {
	repo       UserRepo
	ucKeys     keys.UseCase
	ldapClient ldap.LdapClient
	// tokenMarker token.Maker
}

// SignIn implements UseCase.
func (s *service) SignIn(ctx context.Context, email string, password string) (*sharedkernel.UserAuth, error) {
	slog.Info("Service Auth:: SignIn")
	isEmailCompany := checkEmailCompany(email)
	if !isEmailCompany {
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
		publicKey, privateKey, err := createKeyPair()
		if err != nil {
			return nil, status.Error(codes.Unknown, "Create Public Private Key failed")
		}
		slog.Info("Public Key", publicKey)
		slog.Info("Private Key", privateKey)

	} else {
		auth, _, err := s.ldapClient.Authenticate(email, password)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		if auth {
			foundUser, err := findUserByEmail(s.repo, ctx, email)
			if err != nil {
				return nil, err
			}
			slog.Info("DATA::", foundUser)
			publicKey, privateKey, err := createKeyPair()
			if err != nil {
				return nil, status.Error(codes.Unknown, "Create Public Private Key failed")
			}
			slog.Info("Public Key", publicKey)
			slog.Info("Private Key", privateKey)

		}

	}
	return nil, status.Error(codes.Canceled, "Test API")
}

// SignUp implements UseCase.
func (s *service) SignUp(ctx context.Context, email, password, fistName, lastName string) (*sharedkernel.UserAuth, error) {
	slog.Info("Service Auth:: SignUp")
	isEmailCompany := checkEmailCompany(email)
	if !isEmailCompany {
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
		publicKey, err := utils.GenerateRandomHexBytes(64)
		if err != nil {
			return nil, status.Error(codes.FailedPrecondition, "Create public key failed")
		}
		privateKey, err := utils.GenerateRandomHexBytes(64)
		if err != nil {
			return nil, status.Error(codes.FailedPrecondition, "Create private key failed")
		}
		modelKey := domain.NewKey(newUser.ID, publicKey, privateKey)
		keyStore, err := s.ucKeys.CreateKeyToken(ctx, modelKey)
		if err != nil {
			return nil, status.Error(codes.FailedPrecondition, "Create key token failed")
		}
		slog.Info("Create key token::", keyStore)
		// tokens , err :=
		// return &sharedkernel.UserAuth{
		// 	User:        newUser,
		// 	AccessToken: tokens.,
		// }
	}
	return nil, status.Error(codes.AlreadyExists, fmt.Sprintf("Email has suffix is %s does't created. Please use the login function for this type of email!", suffixEmailCompany))
}

var _ UseCase = (*service)(nil)

func NewUseCase(
	repo UserRepo,
	ucKeys keys.UseCase,
	ldapClient ldap.LdapClient,
	// tokenMarker token.Maker,
) UseCase {
	return &service{
		repo:       repo,
		ucKeys:     ucKeys,
		ldapClient: ldapClient,
		// tokenMarker: tokenMarker,
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

// func createTokenPair(tokenMarker token.Maker, payload interface{}, publicKey, privateKey string) (string, string) {
// 	accessToken := tokenMarker.CreateToken(payload)
// }
