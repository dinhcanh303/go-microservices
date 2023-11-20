package sharedkernel

import (
	"github.com/dinhcanh303/go-microservices/internal/auth/domain"
)

type UserAuth struct {
	User         *domain.User `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
}
type Payload struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FullName  string `json:"full_name"`
	Role      string `json:"role"`
	AvatarUrl string `json:"avatar_url"`
}
