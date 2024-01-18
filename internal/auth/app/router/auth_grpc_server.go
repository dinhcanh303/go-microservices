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
	cfg                 *config.Config
	uc                  auth.UseCase
	ucKey               keys.UseCase
	uploadDomainService domain.UploadDomainService
}

var _ gen.AuthServiceServer = (*authGRPCServer)(nil)

var AuthGRPCServerSet = wire.NewSet(NewAuthGRPCServer)

func NewAuthGRPCServer(
	grpcServer *grpc.Server,
	cfg *config.Config,
	uc auth.UseCase,
	ucKey keys.UseCase,
	uploadDomainService domain.UploadDomainService) gen.AuthServiceServer {
	svc := authGRPCServer{
		cfg:                 cfg,
		uc:                  uc,
		ucKey:               ucKey,
		uploadDomainService: uploadDomainService,
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
	signUpRes, err := a.uc.SignUp(ctx, request.Email, request.Password, request.FirstName, request.LastName)
	if err != nil {
		return nil, err
	}
	return &gen.SignUpResponse{
		User: &gen.User{
			Id:         signUpRes.User.ID.String(),
			Email:      signUpRes.User.Email,
			FirstName:  signUpRes.User.FirstName,
			LastName:   signUpRes.User.LastName,
			FullName:   signUpRes.User.FullName,
			NickName:   signUpRes.User.NickName,
			Role:       signUpRes.User.Role,
			AvatarUrl:  signUpRes.User.AvatarUrl,
			ProfileUrl: signUpRes.User.ProfileUrl,
			CreatedAt:  timestamppb.New(signUpRes.User.CreatedAt),
			UpdatedAt:  timestamppb.New(signUpRes.User.UpdatedAt),
		},
		AccessToken:  signUpRes.AccessToken,
		RefreshToken: signUpRes.RefreshToken,
	}, nil
}
func (a *authGRPCServer) SignIn(ctx context.Context, request *gen.SignInRequest) (*gen.SignInResponse, error) {
	slog.Info("POST:: SignIn")
	violations := validation.ValidateSignIn(request)
	if violations != nil {
		return nil, errorPkg.InvalidArgumentError(violations)
	}
	signInRes, err := a.uc.SignIn(ctx, request.Email, request.Password)
	if err != nil {
		return nil, err
	}
	// avatarUrl, thumbnailUrl := getAvatarAndThumbnailAvatar(a, ctx, signInRes.User.ID)
	return &gen.SignInResponse{
		User: &gen.User{
			Id:         signInRes.User.ID.String(),
			Email:      signInRes.User.Email,
			FirstName:  signInRes.User.FirstName,
			LastName:   signInRes.User.LastName,
			FullName:   signInRes.User.FullName,
			NickName:   signInRes.User.NickName,
			Role:       signInRes.User.Role,
			AvatarUrl:  signInRes.User.AvatarUrl,
			ProfileUrl: signInRes.User.ProfileUrl,
			CreatedAt:  timestamppb.New(signInRes.User.CreatedAt),
			UpdatedAt:  timestamppb.New(signInRes.User.UpdatedAt),
		},
		AccessToken:  signInRes.AccessToken,
		RefreshToken: signInRes.RefreshToken,
	}, nil
}
func (a *authGRPCServer) GetUsers(ctx context.Context, request *gen.GetUsersRequest) (*gen.GetUsersResponse, error) {
	slog.Info("GET:: GetUsers")
	users, err := a.uc.GetUsers(ctx, request.GetSearch(), request.GetLimit(), request.GetOffset())
	if err != nil {
		return nil, err
	}
	return &gen.GetUsersResponse{
		Users: lo.Map(users, func(user *domain.User, _ int) *gen.User {
			// avatarUrl, thumbnailUrl := getAvatarAndThumbnailAvatar(a, ctx, user.ID)
			return &gen.User{
				Id:          user.ID.String(),
				Email:       user.Email,
				FirstName:   user.FirstName,
				LastName:    user.LastName,
				FullName:    user.FullName,
				NickName:    user.NickName,
				Role:        user.Role,
				AvatarUrl:   user.AvatarUrl,
				ProfileUrl:  user.ProfileUrl,
				Gender:      user.Gender,
				Phone:       user.Phone,
				Address:     user.Address,
				DateOfBirth: timestamppb.New(user.DateOfBirth),
				Position:    user.Position,
				CreatedAt:   timestamppb.New(user.CreatedAt),
				UpdatedAt:   timestamppb.New(user.UpdatedAt),
			}
		}),
	}, nil

}

func (a *authGRPCServer) UpdateUser(ctx context.Context, request *gen.UpdateUserRequest) (*gen.UpdateUserResponse, error) {
	slog.Info("GET:: UpdateUser")
	slog.Info("REQUEST::", request)
	userIdReq, err := uuid.Parse(request.User.Id)
	if err != nil {
		return nil, err
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("no headers found in the incoming context.")
	}
	clientId := utils.GetKeyMetadata(md, constant.ClientID)
	if clientId == "" {
		return nil, errors.New("No found client ID")
	}
	userId, err := uuid.Parse(clientId)
	if err != nil {
		return nil, err
	}
	if userId != userIdReq {
		return nil, errors.New("No matched ID please check ID request")
	}
	model := &domain.User{
		ID:          userId,
		AvatarUrl:   request.User.AvatarUrl,
		ProfileUrl:  request.User.ProfileUrl,
		Gender:      request.User.Gender,
		Phone:       request.User.Phone,
		Address:     request.User.Address,
		DateOfBirth: request.User.DateOfBirth.AsTime(),
		Position:    request.User.Position,
	}
	user, err := a.uc.UpdateUser(ctx, model)
	if err != nil {
		return nil, err
	}
	return &gen.UpdateUserResponse{
		User: &gen.User{
			Id:          user.ID.String(),
			Email:       user.Email,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			FullName:    user.FullName,
			NickName:    user.NickName,
			Role:        user.Role,
			AvatarUrl:   user.AvatarUrl,
			ProfileUrl:  user.ProfileUrl,
			Gender:      user.Gender,
			Phone:       user.Phone,
			Address:     user.Address,
			DateOfBirth: timestamppb.New(user.DateOfBirth),
			Position:    user.Position,
			CreatedAt:   timestamppb.New(user.CreatedAt),
			UpdatedAt:   timestamppb.New(user.UpdatedAt),
		},
	}, nil
}

func (a *authGRPCServer) GetProfile(ctx context.Context, request *gen.GetProfileRequest) (*gen.GetProfileResponse, error) {
	slog.Info("GET:: Profile")
	var userId uuid.UUID
	var err error
	if request.Id != "" {
		userId, err = uuid.Parse(request.Id)
		if err != nil {
			return nil, err
		}
	} else {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.New("no headers found in the incoming context.")
		}
		clientId := utils.GetKeyMetadata(md, constant.ClientID)
		if clientId == "" {
			return nil, errors.New("Invalid Request")
		}
		userId, err = uuid.Parse(clientId)
		if err != nil {
			return nil, err
		}
	}

	user, err := a.uc.GetUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	// avatarUrl, thumbnailUrl := getAvatarAndThumbnailAvatar(a, ctx, userId)
	return &gen.GetProfileResponse{
		User: &gen.User{
			Id:          user.ID.String(),
			Email:       user.Email,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			FullName:    user.FullName,
			NickName:    user.NickName,
			Role:        user.Role,
			AvatarUrl:   user.AvatarUrl,
			ProfileUrl:  user.ProfileUrl,
			Gender:      user.Gender,
			Phone:       user.Phone,
			Address:     user.Address,
			DateOfBirth: timestamppb.New(user.DateOfBirth),
			Position:    user.Position,
			CreatedAt:   timestamppb.New(user.CreatedAt),
			UpdatedAt:   timestamppb.New(user.UpdatedAt),
		},
	}, nil
}
func (a *authGRPCServer) Verify(ctx context.Context, request *gen.VerifyRequest) (*gen.VerifyResponse, error) {
	slog.Info("GET:: Verify")
	md, ok := metadata.FromIncomingContext(ctx)
	slog.Info("Payload::", md)
	if !ok {
		return nil, errors.New("no headers found in the incoming context.")
	}
	clientId := utils.GetKeyMetadata(md, constant.ClientID)
	if clientId == "" {
		return nil, errors.New("Invalid Request")
	}
	slog.Info("clientId::", clientId)
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
		payload, err := verifyToken(refreshToken, keyStore.PrivateKey)
		if err != nil {
			return nil, errors.New("Unauthorized")
		}
		grpc.SendHeader(ctx, addHeader(payload, keyStore, refreshToken))
		return &gen.VerifyResponse{}, nil
	}
	authorization := utils.GetKeyMetadata(md, constant.Authorization)
	slog.Info("authorization::", authorization)
	if authorization == "" {
		return nil, errors.New("Unauthorized")
	}
	payload, err := verifyToken(authorization, keyStore.PublicKey)
	if err != nil {
		return nil, errors.New("Unauthorized")
	}
	slog.Info("Payload::", payload)
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
			Id:          res.User.ID.String(),
			Email:       res.User.Email,
			FirstName:   res.User.FirstName,
			LastName:    res.User.LastName,
			FullName:    res.User.FullName,
			NickName:    res.User.NickName,
			Role:        res.User.Role,
			AvatarUrl:   res.User.AvatarUrl,
			ProfileUrl:  res.User.ProfileUrl,
			Gender:      res.User.Gender,
			Phone:       res.User.Phone,
			Address:     res.User.Address,
			DateOfBirth: timestamppb.New(res.User.DateOfBirth),
			Position:    res.User.Position,
			CreatedAt:   timestamppb.New(res.User.CreatedAt),
			UpdatedAt:   timestamppb.New(res.User.UpdatedAt),
		},
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}, nil
}
func (a *authGRPCServer) GetUserIdsOfCompanyByUserId(ctx context.Context, request *gen.GetUserIdsOfCompanyByUserIdRequest) (*gen.GetUserIdsOfCompanyByUserIdResponse, error) {
	slog.Info("GET:: GetUserIdsOfCompanyByUserId")
	userId, err := uuid.Parse(request.UserId)
	if err != nil {
		return nil, errors.New("failed to parse uuid")
	}
	userIds, err := a.uc.GetUserIdsOfCompanyByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	return &gen.GetUserIdsOfCompanyByUserIdResponse{
		UserIds: lo.Map(userIds, func(item uuid.UUID, _ int) string {
			return item.String()
		}),
	}, nil
}

// func getAvatarAndThumbnailAvatar(a *authGRPCServer, ctx context.Context, userId uuid.UUID) (string, string) {
// 	avatarRes, err := a.uploadDomainService.GetAvatarUser(ctx, userId)
// 	if err != nil {
// 		slog.Warn("uploadDomainService.GetAvatarUser failed", err)
// 		return "", ""
// 	}
// 	return avatarRes.URL, avatarRes.URLThumbnail
// }
