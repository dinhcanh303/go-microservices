package router

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/dinhcanh303/go-microservices/cmd/gateway/config"
	"github.com/dinhcanh303/go-microservices/internal/auth/app/validation"
	"github.com/dinhcanh303/go-microservices/internal/auth/usecases/auth"
	"github.com/dinhcanh303/go-microservices/internal/auth/usecases/keys"
	errorPkg "github.com/dinhcanh303/go-microservices/pkg/error"
	"github.com/dinhcanh303/go-microservices/proto/gen"
	"github.com/google/uuid"
	"github.com/google/wire"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type authGRPCServer struct {
	gen.UnimplementedAuthServiceServer
	cfg   *config.Config
	uc    auth.UseCase
	ucKey keys.UseCase
}

var _ gen.AuthServiceServer = (*authGRPCServer)(nil)

var AuthGRPCServerSet = wire.NewSet(NewAuthGRPCServer)

func NewAuthGRPCServer(
	grpcServer *grpc.Server,
	cfg *config.Config,
	uc auth.UseCase,
	ucKey keys.UseCase) gen.AuthServiceServer {
	svc := authGRPCServer{
		cfg:   cfg,
		uc:    uc,
		ucKey: ucKey,
	}
	gen.RegisterAuthServiceServer(grpcServer, &svc)
	reflection.Register(grpcServer)
	return &svc
}

func (a *authGRPCServer) SignUp(ctx context.Context, request *gen.SignUpRequest) (*gen.SignUpResponse, error) {
	slog.Info("POST:: SignUp")
	// violations := validation.ValidateSignUp(request)
	// if violations != nil {
	// 	return nil, errorPkg.InvalidArgumentError(violations)
	// }
	signRes, err := a.uc.SignUp(ctx, request.Email, request.Password, request.FirstName, request.LastName)
	if err != nil {
		return nil, err
	}
	return &gen.SignUpResponse{
		User: &gen.User{
			Id:        signRes.User.ID.String(),
			Email:     signRes.User.Email,
			FirstName: signRes.User.FirstName,
			LastName:  signRes.User.LastName,
			FullName:  signRes.User.FullName,
			CreatedAt: timestamppb.New(signRes.User.CreatedAt),
			UpdatedAt: timestamppb.New(signRes.User.UpdatedAt),
		},
		AccessToken:  signRes.AccessToken,
		RefreshToken: signRes.RefreshToken,
	}, nil
}
func (a *authGRPCServer) SignIn(ctx context.Context, request *gen.SignInRequest) (*gen.SignInResponse, error) {
	slog.Info("POST:: SignIn")
	violations := validation.ValidateSignIn(request)
	if violations != nil {
		return nil, errorPkg.InvalidArgumentError(violations)
	}
	signRes, err := a.uc.SignIn(ctx, request.Email, request.Password)
	if err != nil {
		return nil, err
	}
	return &gen.SignInResponse{
		User: &gen.User{
			Id:        signRes.User.ID.String(),
			Email:     signRes.User.Email,
			FirstName: signRes.User.FirstName,
			LastName:  signRes.User.LastName,
			FullName:  signRes.User.FullName,
			CreatedAt: timestamppb.New(signRes.User.CreatedAt),
			UpdatedAt: timestamppb.New(signRes.User.UpdatedAt),
		},
		AccessToken:  signRes.AccessToken,
		RefreshToken: signRes.RefreshToken,
	}, nil
}

// func (a *authGRPCServer) HandleRefreshToken(ctx context.Context) {
// 	userId, email := "", ""
// 	if slices.Contains(keyStore) {

//		}
//	}
func (a *authGRPCServer) FindKeyByUserID(ctx context.Context, request *gen.FindKeyByUserIDRequest) (*gen.FindKeyByUserIDResponse, error) {
	slog.Info("GET:: FindKeyByUserID")
	userId, err := uuid.Parse(request.UserId)
	if err != nil {
		return nil, err
	}
	keyStore, err := a.ucKey.FindKeyByUserID(ctx, userId)
	if err != nil {
		return nil, err
	}
	var stringSlice []string
	err = json.Unmarshal(keyStore.RefreshTokensUsed, &stringSlice)
	if err != nil {
		return nil, err
	}
	return &gen.FindKeyByUserIDResponse{
		Id:               int32(keyStore.ID),
		UserId:           keyStore.UserID.String(),
		PublicKey:        keyStore.PublicKey,
		PrivateKey:       keyStore.PrivateKey,
		RefreshToken:     keyStore.RefreshToken,
		RefreshTokenUsed: stringSlice,
		CreatedAt:        timestamppb.New(keyStore.CreatedAt),
		UpdatedAt:        timestamppb.New(keyStore.UpdatedAt),
	}, nil
}
