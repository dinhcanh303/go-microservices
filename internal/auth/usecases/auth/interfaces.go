package auth

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/auth/domain"
	"github.com/dinhcanh303/go-microservices/pkg/rabbitmq/publisher"
	"github.com/google/uuid"
)

type UserRepo interface {
	CreateUser(context.Context, *domain.User) (*domain.User, error)
	GetUser(ctx context.Context, userId uuid.UUID) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetAllUserIdOfCompany(ctx context.Context, company string) ([]uuid.UUID, error)
}
type UseCase interface {
	SignIn(ctx context.Context, email, password string) (*domain.UserAuth, error)
	SignUp(ctx context.Context, email, password, fistName, lastName string) (*domain.UserAuth, error)
	HandleRefreshToken(ctx context.Context, email, refreshToken string) (*domain.UserAuth, error)
	Logout(ctx context.Context, key *domain.Key) error
	GetAllUserIdOfCompanyByUserId(ctx context.Context, userId uuid.UUID) ([]uuid.UUID, error)
	GetUser(ctx context.Context, userID uuid.UUID) (*domain.User, error)
}

type UserCreatedEventPublisher interface {
	Configure(...publisher.Option)
	Publish(context.Context, []byte, string) error
}
type UserDeletedEventPublisher interface {
	Configure(...publisher.Option)
	Publish(context.Context, []byte, string) error
}
