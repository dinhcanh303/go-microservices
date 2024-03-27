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

var UploadHandlerSet = wire.NewSet(NewAuthHandler)

func NewAuthHandler(uc auth.UseCaseHttp) auth.UseCaseHttp {
	return &AuthHandler{
		uc: uc,
	}
}

// GoogleCallback implements auth.UseCaseHttp.
func (a *AuthHandler) GoogleCallback(res http.ResponseWriter, req *http.Request) {
	panic("unimplemented")
}

// GoogleLogin implements auth.UseCaseHttp.
func (a *AuthHandler) GoogleLogin(res http.ResponseWriter, req *http.Request) {

	panic("unimplemented")
	// http.Redirect(res, req, url, http.StatusSeeOther)
}
