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
	UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUsers(ctx context.Context, search string, limit, offset int32) ([]*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserIdsOfCompany(ctx context.Context, company string) ([]uuid.UUID, error)
	GetUsersInviteGroup(ctx context.Context, groupIds []uuid.UUID, limit, offset int32) ([]*domain.User, error)
	GetUsersBirthDayByCurrentMonth(ctx context.Context) ([]*domain.User, error)
	GetUsersBirthDayByCurrentDay(ctx context.Context) ([]*domain.User, error)
}
type UseCase interface {
	SignIn(ctx context.Context, email, password string) (*domain.UserAuth, error)
	SignUp(ctx context.Context, email, password, fistName, lastName string) (*domain.UserAuth, error)
	HandleRefreshToken(ctx context.Context, email, refreshToken string) (*domain.UserAuth, error)
	Logout(ctx context.Context, key *domain.Key) error
	GetUserIdsOfCompanyByUserId(ctx context.Context, userId uuid.UUID) ([]uuid.UUID, error)
	GetUser(ctx context.Context, userID uuid.UUID) (*domain.User, error)
	GetUsers(ctx context.Context, search string, limit, offset int32) ([]*domain.User, error)
	GetUsersInviteGroup(ctx context.Context, groupIds []uuid.UUID, limit, offset int32) ([]*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUsersBirthDayByCurrentMonth(ctx context.Context) ([]*domain.User, error)
	GetUsersBirthDayByCurrentDay(ctx context.Context) ([]*domain.User, error)
}

type UserCreatedEventPublisher interface {
	Configure(...publisher.Option)
	Publish(context.Context, []byte, string) error
}
type UserDeletedEventPublisher interface {
	Configure(...publisher.Option)
	Publish(context.Context, []byte, string) error
}
