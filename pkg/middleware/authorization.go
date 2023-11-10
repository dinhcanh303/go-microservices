package middleware

import (
	"net/http"

	"github.com/dinhcanh303/go-microservices/pkg/constant"
)

func AuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := r.Header.Get(constant.ClientID)
		if userId == "" {
			http.Error(w, "Invalid client", 200)
			return
		}
		h.ServeHTTP(w, r)
	})
}
