package router

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/dinhcanh303/go-microservices/cmd/auth/config"
	"github.com/dinhcanh303/go-microservices/internal/auth/app/validation"
	"github.com/dinhcanh303/go-microservices/internal/auth/domain"
	"github.com/dinhcanh303/go-microservices/internal/auth/usecases/auth"
	"github.com/dinhcanh303/go-microservices/internal/auth/usecases/keys"
	"github.com/dinhcanh303/go-microservices/pkg/constant"
	errorPkg "github.com/dinhcanh303/go-microservices/pkg/error"
	"github.com/dinhcanh303/go-microservices/pkg/token"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
	"github.com/dinhcanh303/go-microservices/proto/gen"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/samber/lo"
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
	slog.Info("POST:: SignIn 2")
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
	clientId := utils.GetKeyMetadata(md, constant.ClientID)
	if clientId == "" {
		return nil, errors.New("Invalid Request")
	}
	userId, err := uuid.Parse(clientId)
	if err != nil {
		return nil, err
	}
	keyStore, err := a.ucKey.FindKeyByUserID(ctx, userId)
	if err != nil {
		return nil, errors.Wrap(err, "keystore::")
	}
	refreshToken := utils.GetKeyMetadata(md, constant.RefreshToken)
	if refreshToken != "" {
		slog.Info("refresh token::", refreshToken)
		payload, err := verifyToken(refreshToken, keyStore.PrivateKey)
		if err != nil {
			return nil, errors.New("Unauthorized")
		}
		grpc.SendHeader(ctx, addHeader(payload, keyStore, refreshToken))
		return &gen.VerifyResponse{}, nil
	}
	authorization := utils.GetKeyMetadata(md, constant.Authorization)
	if authorization == "" {
		return nil, errors.New("Unauthorized")
	}
	payload, err := verifyToken(authorization, keyStore.PublicKey)
	if err != nil {
		return nil, errors.New("Unauthorized")
	}
	grpc.SendHeader(ctx, addHeader(payload, keyStore, ""))
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
func addHeader(payload *token.Payload, keyStore *domain.Key, refreshToken string) metadata.MD {
	header := metadata.Pairs(
		constant.User,
		fmt.Sprintf("%s,%s,%s,%s,%s",
			payload.ID.String(), payload.Email,
			payload.FullName, payload.Role,
			payload.AvatarUrl),
		constant.KeyStore,
		fmt.Sprintf("%s,%s,%s,%s,%s,%s",
			strconv.FormatInt(keyStore.ID, 10), keyStore.UserID,
			keyStore.PublicKey, keyStore.PrivateKey,
			keyStore.RefreshToken, keyStore.RefreshTokensUsed,
		),
		constant.RefreshToken,
		refreshToken,
	)
	return header
}
func (a *authGRPCServer) Logout(ctx context.Context, request *gen.LogoutRequest) (*gen.LogoutResponse, error) {
	slog.Info("GET:: Logout")
	keyStore, err := utils.ExtractMetadataKeyStore(ctx)
	if err != nil {
		return nil, err
	}
	slog.Info("KeyStore::", keyStore)
	err = a.ucKey.DeleteKeyByID(ctx, keyStore.ID)
	if err != nil {
		return nil, errors.Wrap(err, "Logout failed :")
	}
	return &gen.LogoutResponse{}, nil
}
func (a *authGRPCServer) HandleRefreshToken(ctx context.Context, request *gen.HandleRefreshTokenRequest) (*gen.HandleRefreshTokenResponse, error) {
	slog.Info("GET:: HandleRefreshToken")
	user, err := utils.ExtractMetadataUser(ctx)
	if err != nil {
		return nil, err
	}
	keyStore, err := utils.ExtractMetadataKeyStore(ctx)
	if err != nil {
		return nil, err
	}
	refreshToken, err := utils.ExtractMetadataRefreshToken(ctx)
	if err != nil {
		return nil, err
	}
	if lo.Contains(keyStore.RefreshTokensUsed, refreshToken) {
		err := a.ucKey.DeleteKeyByUserID(ctx, user.ID)
		if err != nil {
			return nil, err
		}
		return nil, errors.New("Something wrong happened !! Please login again")
	}
	if keyStore.RefreshToken != refreshToken {
		return nil, errors.New("User not registered")
	}

	slog.Info("User::", user)
	slog.Info("KeyStore::", keyStore.RefreshTokensUsed)
	res, err := a.uc.HandleRefreshToken(ctx, user.Email, refreshToken)
	if err != nil {
		return nil, err
	}
	return &gen.HandleRefreshTokenResponse{
		User: &gen.User{
			Id:        res.User.ID.String(),
			Email:     res.User.Email,
			FirstName: res.User.FirstName,
			LastName:  res.User.LastName,
			FullName:  res.User.FullName,
			CreatedAt: timestamppb.New(res.User.CreatedAt),
			UpdatedAt: timestamppb.New(res.User.UpdatedAt),
		},
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}, nil
}
func (a *authGRPCServer) GetAllUserIdOfCompanyByUserId(ctx context.Context, request *gen.GetAllUserIdOfCompanyByUserIdRequest) (*gen.GetAllUserIdOfCompanyByUserIdResponse, error) {
	slog.Info("GET:: GetAllUserIdByUserId")
	userId, err := uuid.Parse(request.UserId)
	if err != nil {
		return nil, errors.New("failed to parse uuid")
	}
	userIds, err := a.uc.GetAllUserIdOfCompanyByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	return &gen.GetAllUserIdOfCompanyByUserIdResponse{
		UserIds: lo.Map(userIds, func(item uuid.UUID, _ int) string {
			return item.String()
		}),
	}, nil
}
