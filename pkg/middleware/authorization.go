package middleware

import (
	"context"
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
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		refreshToken := r.Header.Get(constant.RefreshToken)
		if refreshToken != "" {
			payload, err := verifyToken(refreshToken, keyStore.PrivateKey)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			ctx := setValueContext(r, payload, keyStore, refreshToken)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		authorization := r.Header.Get(constant.Authorization)
		if authorization == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		payload, err := verifyToken(authorization, keyStore.PublicKey)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := setValueContext(r, payload, keyStore, "")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func setValueContext(r *http.Request, payload *token.Payload, key *gen.FindKeyByUserIDResponse, refreshToken string) context.Context {
	ctx := context.WithValue(r.Context(), constant.KeyStore, key)
	ctx = context.WithValue(ctx, constant.User, payload)
	if refreshToken != "" {
		ctx = context.WithValue(ctx, constant.RefreshToken, refreshToken)
	}
	return ctx
}
func verifyToken(refreshToken, secretKey string) (*token.Payload, error) {
	jwt := token.NewJWTMaker()
	payload, err := jwt.VerifyToken(refreshToken, secretKey)
	if err != nil {
		return nil, err
	}
	return payload, nil
}
