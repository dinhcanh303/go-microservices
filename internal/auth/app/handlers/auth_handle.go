package handlers

import (
	"net/http"

	"github.com/dinhcanh303/go-microservices/internal/auth/usecases/auth"
	"github.com/google/wire"
)

type AuthHandler struct {
	uc auth.UseCaseHttp
}

var _ auth.UseCaseHttp = (*AuthHandler)(nil)

var AuthHandlerSet = wire.NewSet(NewAuthHandler)

func NewAuthHandler(uc auth.UseCaseHttp) *AuthHandler {
	return &AuthHandler{
		uc: uc,
	}
}

// GoogleCallback implements auth.UseCaseHttp.
func (a *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	a.uc.GoogleCallback(w, r)
}

// GoogleLogin implements auth.UseCaseHttp.
func (a *AuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	a.uc.GoogleLogin(w, r)
}
