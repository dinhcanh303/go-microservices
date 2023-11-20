package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/dinhcanh303/go-microservices/internal/auth/app"
	"github.com/dinhcanh303/go-microservices/pkg/constant"
	"github.com/dinhcanh303/go-microservices/pkg/token"
	"github.com/dinhcanh303/go-microservices/proto/gen"
)

func AuthMiddleware(next http.Handler, app *app.App) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := r.Header.Get(constant.ClientID)
		if userId == "" {
			http.Error(w, "Invalid client ID", http.StatusUnauthorized)
			return
		}
		keyStore, err := app.AuthGRPCServer.FindKeyByUserID(context.Background(), &gen.FindKeyByUserIDRequest{
			UserId: userId,
		})
		if err != nil {
			http.Error(w, "Not Found KeyStore", http.StatusNotFound)
		}
		refreshToken := r.Header.Get(constant.RefreshToken)
		if refreshToken != "" {
			payload, err := verifyToken(refreshToken, keyStore.PrivateKey)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			payloadByte, err := json.Marshal(payload)
			r.Header.Set("keyStore", string(payloadByte))
		}
		authorization := r.Header.Get(constant.Authorization)
		if authorization == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
func verifyToken(refreshToken, secretKey string) (*token.Payload, error) {
	jwt := token.NewJWTMaker()
	payload, err := jwt.VerifyToken(refreshToken, secretKey)
	if err != nil {
		return nil, err
	}
	return payload, nil
}
