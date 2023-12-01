package router

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/dinhcanh303/go-microservices/cmd/auth/config"
	"github.com/dinhcanh303/go-microservices/internal/auth/app/validation"
	"github.com/dinhcanh303/go-microservices/internal/auth/usecases/auth"
	"github.com/dinhcanh303/go-microservices/internal/auth/usecases/keys"
	"github.com/dinhcanh303/go-microservices/pkg/constant"
	errorPkg "github.com/dinhcanh303/go-microservices/pkg/error"
	"github.com/dinhcanh303/go-microservices/pkg/token"
	"github.com/dinhcanh303/go-microservices/proto/gen"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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
func (a *authGRPCServer) Verify(ctx context.Context, request *gen.VerifyRequest) (*gen.VerifyResponse, error) {
	slog.Info("GET:: Verify")
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("no headers found in the incoming context.")
	}
	slog.Info("Metadata::", md)
	clientId := ""
	if values := md.Get(constant.ClientID); len(values) > 0 {
		clientId = values[0]
	}
	if clientId == "" {
		return nil, errors.New("Invalid Request")
	}
	userId, err := uuid.Parse(clientId)
	if err != nil {
		return nil, err
	}
	keyStore, err := a.ucKey.FindKeyByUserID(ctx, userId)
	if err != nil {
		return nil, err
	}
	refreshToken := ""
	if values := md.Get(constant.RefreshToken); len(values) > 0 {
		refreshToken = values[0]
	}
	if refreshToken != "" {
		payload, err := verifyToken(refreshToken, keyStore.PrivateKey)
		if err != nil {
			return nil, errors.New("Unauthorized")
		}
		grpc.SendHeader(ctx, addHeader(payload))
		return &gen.VerifyResponse{}, nil
	}
	authorization := ""
	if values := md.Get(constant.Authorization); len(values) > 0 {
		authorization = values[0]
	}
	if authorization == "" {
		return nil, errors.New("Unauthorized")
	}
	payload, err := verifyToken(authorization, keyStore.PublicKey)
	if err != nil {
		return nil, errors.New("Unauthorized")
	}
	slog.Info("TEST RESPONSE HEADER:")
	metadata.NewOutgoingContext(ctx, addHeader(payload))
	grpc.SendHeader(ctx, addHeader(payload))
	return &gen.VerifyResponse{}, nil
}
func verifyToken(refreshToken, secretKey string) (*token.Payload, error) {
	jwt := token.NewJWTMaker()
	payload, err := jwt.VerifyToken(refreshToken, secretKey)
	if err != nil {
		return nil, err
	}
	return payload, nil
}
func addHeader(payload *token.Payload) metadata.MD {
	header := metadata.Pairs(
		constant.User,
		fmt.Sprintf("%s,%s,%s,%s,%s",
			payload.ID.String(), payload.Email,
			payload.FullName, payload.Role,
			payload.AvatarUrl),
	)
	return header
}
