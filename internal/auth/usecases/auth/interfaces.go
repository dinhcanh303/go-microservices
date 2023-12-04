package auth

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/auth/domain"
	sharedkernel "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
	"github.com/google/uuid"
)

type UserRepo interface {
	CreateUser(context.Context, *domain.User) (*domain.User, error)
	GetUser(ctx context.Context, userID uuid.UUID) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetAllUserIdOfCompany(ctx context.Context, company string) ([]uuid.UUID, error)
}
type UseCase interface {
	SignIn(ctx context.Context, email, password string) (*sharedkernel.UserAuth, error)
	SignUp(ctx context.Context, email, password, fistName, lastName string) (*sharedkernel.UserAuth, error)
	Logout(ctx context.Context, key *domain.Key) error
	GetAllUserIdByUserId(ctx context.Context, userId uuid.UUID) ([]uuid.UUID, error)
}
