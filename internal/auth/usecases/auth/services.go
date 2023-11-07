package auth

import (
	"context"
	"strings"

	"github.com/dinhcanh303/go-microservices/internal/auth/domain"
	"github.com/dinhcanh303/go-microservices/internal/auth/usecases/keys"
	sharedkernel "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
	"github.com/dinhcanh303/go-microservices/pkg/ldap"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type service struct {
	repo       UserRepo
	ucKeys     keys.UseCase
	ldapClient ldap.LdapClient
}

// SignIn implements UseCase.
func (s *service) SignIn(ctx context.Context, email string, password string) (*sharedkernel.UserAuth, error) {
	isEmailCompany := checkEmailCompany(email)
	if !isEmailCompany {
		foundUser, err := s.repo.GetUserByEmail(ctx, email)
		if err != nil {
			return nil, status.Error(codes.NotFound, "User not registered")
		}

	} else {
		auth, info, err := s.ldapClient.Authenticate(email, password)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
	}

}

// SignUp implements UseCase.
func (s *service) SignUp(ctx context.Context, email, password, fistName, lastName string) (*sharedkernel.UserAuth, error) {
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
		_, err = s.ucKeys.CreateKeyToken(ctx, modelKey)
		if err != nil {
			return nil, status.Error(codes.FailedPrecondition, "Create key token failed")
		}
		tokens := 

	}
	panic("unimplemented")
}

var _ UseCase = (*service)(nil)

func checkEmailCompany(email string) bool {
	return strings.Contains(email, "@tlcmodular.com")
}
