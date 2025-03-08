package router

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"

	v1 "github.com/dinhcanh303/go-microservices/api/auth/v1"
	v1g "github.com/dinhcanh303/go-microservices/api/group/v1"
	"github.com/dinhcanh303/go-microservices/cmd/auth/config"
	"github.com/dinhcanh303/go-microservices/internal/auth/app/validation"
	"github.com/dinhcanh303/go-microservices/internal/auth/domain"
	"github.com/dinhcanh303/go-microservices/internal/auth/usecases/auth"
	"github.com/dinhcanh303/go-microservices/internal/auth/usecases/follow"
	"github.com/dinhcanh303/go-microservices/internal/auth/usecases/keys"
	"github.com/dinhcanh303/go-microservices/pkg/constant"
	errorPkg "github.com/dinhcanh303/go-microservices/pkg/error"
	"github.com/dinhcanh303/go-microservices/pkg/redis"
	"github.com/dinhcanh303/go-microservices/pkg/token"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type authGRPCServer struct {
	v1.UnimplementedAuthServiceServer
	cfg                 *config.Config
	uc                  auth.UseCase
	ucKey               keys.UseCase
	ucFollow            follow.UseCase
	uploadDomainService domain.UploadDomainService
	groupDomainService  domain.GroupDomainService
	redis               redis.RedisEngine
}

var _ v1.AuthServiceServer = (*authGRPCServer)(nil)

var AuthGRPCServerSet = wire.NewSet(NewAuthGRPCServer)

func NewAuthGRPCServer(
	grpcServer *grpc.Server,
	cfg *config.Config,
	uc auth.UseCase,
	ucKey keys.UseCase,
	ucFollow follow.UseCase,
	uploadDomainService domain.UploadDomainService,
	groupDomainService domain.GroupDomainService,
	redis redis.RedisEngine) v1.AuthServiceServer {
	svc := authGRPCServer{
		cfg:                 cfg,
		uc:                  uc,
		ucKey:               ucKey,
		ucFollow:            ucFollow,
		uploadDomainService: uploadDomainService,
		groupDomainService:  groupDomainService,
		redis:               redis,
	}
	v1.RegisterAuthServiceServer(grpcServer, &svc)
	reflection.Register(grpcServer)
	return &svc
}

func (a *authGRPCServer) SignUp(ctx context.Context, request *v1.SignUpRequest) (*v1.SignUpResponse, error) {
	slog.Info("POST:: SignUp")
	// violations := validation.ValidateSignUp(request)
	// if violations != nil {
	// 	return nil, errorPkg.InvalidArgumentError(violations)
	// }
	signUpRes, err := a.uc.SignUp(ctx, request.Email, request.Password, request.FirstName, request.LastName)
	if err != nil {
		return nil, err
	}
	return &v1.SignUpResponse{
		User: &v1.User{
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
func (a *authGRPCServer) SignIn(ctx context.Context, request *v1.SignInRequest) (*v1.SignInResponse, error) {
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
	return &v1.SignInResponse{
		User: &v1.User{
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
func (a *authGRPCServer) GetUsers(ctx context.Context, request *v1.GetUsersRequest) (*v1.GetUsersResponse, error) {
	slog.Info("GET:: GetUsers")
	users, err := a.uc.GetUsers(ctx, request.GetSearch(), request.GetLimit(), request.GetOffset())
	if err != nil {
		return nil, errors.Wrap(err, "uc.GetUsers failed")
	}
	return &v1.GetUsersResponse{
		Users: lo.Map(users, func(user *domain.User, _ int) *v1.User {
			// avatarUrl, thumbnailUrl := getAvatarAndThumbnailAvatar(a, ctx, user.ID)
			return &v1.User{
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

func (a *authGRPCServer) UpdateUser(ctx context.Context, request *v1.UpdateUserRequest) (*v1.UpdateUserResponse, error) {
	slog.Info("GET:: UpdateUser")
	userIdReq, err := uuid.Parse(request.User.Id)
	if err != nil {
		return nil, err
	}
	payloadUser, err := utils.ExtractMetadataUser(ctx)
	if err != nil {
		return nil, err
	}
	if payloadUser.ID != userIdReq {
		return nil, errors.New("No matched ID please check ID request")
	}
	model := &domain.User{
		ID:          payloadUser.ID,
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

	return &v1.UpdateUserResponse{
		User: &v1.User{
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

func (a *authGRPCServer) UpdateUserSettings(ctx context.Context, request *v1.UpdateUserSettingsRequest) (*v1.UpdateUserSettingsResponse, error) {
	slog.Info("GET:: UpdateUserSettings")
	payloadUser, err := utils.ExtractMetadataUser(ctx)
	if err != nil {
		return nil, err
	}
	user, err := a.uc.GetUser(ctx, payloadUser.ID)
	if err != nil {
		return nil, err
	}
	var settings domain.Settings
	_ = json.Unmarshal(user.Settings, &settings)
	// if request.Theme != nil {
	// 	settings.Social.System.Theme = request.Theme
	// }
	// if request.StatusPost != nil {
	// 	settings.Social.Post.StatusDefault = request.StatusPost
	// }
	settingJson, _ := json.Marshal(settings)
	model := &domain.User{
		ID:       payloadUser.ID,
		Settings: settingJson,
	}
	userUpdated, err := a.uc.UpdateUser(ctx, model)
	if err != nil {
		return nil, err
	}
	return &v1.UpdateUserSettingsResponse{
		User: &v1.User{
			Id:          userUpdated.ID.String(),
			Email:       userUpdated.Email,
			FirstName:   userUpdated.FirstName,
			LastName:    userUpdated.LastName,
			FullName:    userUpdated.FullName,
			NickName:    userUpdated.NickName,
			Role:        userUpdated.Role,
			AvatarUrl:   userUpdated.AvatarUrl,
			ProfileUrl:  userUpdated.ProfileUrl,
			Gender:      userUpdated.Gender,
			Phone:       userUpdated.Phone,
			Address:     userUpdated.Address,
			DateOfBirth: timestamppb.New(userUpdated.DateOfBirth),
			Position:    userUpdated.Position,
			CreatedAt:   timestamppb.New(userUpdated.CreatedAt),
			UpdatedAt:   timestamppb.New(userUpdated.UpdatedAt),
		},
	}, nil
}
func (a *authGRPCServer) GetUsersBirthDayByCurrentDay(ctx context.Context, request *v1.GetUsersBirthDayByCurrentDayRequest) (*v1.GetUsersBirthDayByCurrentDayResponse, error) {
	users, err := a.uc.GetUsersBirthDayByCurrentDay(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetUsersBirthDayByCurrentDayResponse{
		Users: lo.Map(users, func(user *domain.User, _ int) *v1.User {
			return &v1.User{
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

func (a *authGRPCServer) GetUsersBirthDayByCurrentMonth(ctx context.Context, request *v1.GetUsersBirthDayByCurrentMonthRequest) (*v1.GetUsersBirthDayByCurrentMonthResponse, error) {
	users, err := a.uc.GetUsersBirthDayByCurrentMonth(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.GetUsersBirthDayByCurrentMonthResponse{
		Users: lo.Map(users, func(user *domain.User, _ int) *v1.User {
			return &v1.User{
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
func (a *authGRPCServer) GetUsersInviteByGroupId(ctx context.Context, request *v1.GetUsersInviteGroupIdRequest) (*v1.GetUsersInviteGroupIdResponse, error) {
	groupMembers, err := a.groupDomainService.GetGroupMembers(ctx, request.GroupId)
	if err != nil {
		return &v1.GetUsersInviteGroupIdResponse{}, nil
	}
	groupMemberIdsString := lo.Map(groupMembers.GroupMembers, func(groupMember *v1g.GroupMemberMetadata, index int) string {
		return groupMember.UserId
	})
	groupMemberIds, err := utils.ConvertArStringToArUUID(groupMemberIdsString)
	if err != nil {
		return nil, err
	}
	inviteMembers, err := a.uc.GetUsersInviteGroup(ctx, groupMemberIds, request.Limit, request.Offset)
	if err != nil {
		return nil, err
	}
	return &v1.GetUsersInviteGroupIdResponse{
		Users: lo.Map(inviteMembers, func(user *domain.User, _ int) *v1.User {
			return &v1.User{
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
func (a *authGRPCServer) GetProfile(ctx context.Context, request *v1.GetProfileRequest) (*v1.GetProfileResponse, error) {
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
	settingsAny := &anypb.Any{}
	if user.Settings != nil {
		settingsAny.Value = user.Settings
		settingsAny.TypeUrl = "json.RawMessage"
	}
	return &v1.GetProfileResponse{
		User: &v1.User{
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
			// Settings:    settingsAny,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		},
	}, nil
}
func (a *authGRPCServer) Verify(ctx context.Context, request *v1.VerifyRequest) (*v1.VerifyResponse, error) {
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
		payload, err := verifyToken(refreshToken, keyStore.PrivateKey)
		if err != nil {
			return nil, errors.New("Unauthorized")
		}
		grpc.SendHeader(ctx, addHeader(payload, keyStore, refreshToken))
		return &v1.VerifyResponse{}, nil
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
	return &v1.VerifyResponse{}, nil
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
func (a *authGRPCServer) Logout(ctx context.Context, request *v1.LogoutRequest) (*v1.LogoutResponse, error) {
	slog.Info("GET:: Logout")
	keyStore, err := utils.ExtractMetadataKeyStore(ctx)
	if err != nil {
		return nil, err
	}
	err = a.ucKey.DeleteKeyByID(ctx, keyStore.ID)
	if err != nil {
		return nil, errors.Wrap(err, "Logout failed :")
	}
	return &v1.LogoutResponse{}, nil
}
func (a *authGRPCServer) HandleRefreshToken(ctx context.Context, request *v1.HandleRefreshTokenRequest) (*v1.HandleRefreshTokenResponse, error) {
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

	res, err := a.uc.HandleRefreshToken(ctx, user.Email, refreshToken)
	if err != nil {
		return nil, err
	}
	return &v1.HandleRefreshTokenResponse{
		User: &v1.User{
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

func (a *authGRPCServer) Follow(ctx context.Context, request *v1.FollowRequest) (*v1.FollowResponse, error) {
	slog.Info("POST:: Follow")
	payloadUser, err := utils.ExtractMetadataUser(ctx)
	if err != nil {
		return nil, err
	}
	followingId, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, err
	}
	err = a.ucFollow.Follow(ctx, payloadUser.ID, followingId)
	if err != nil {
		return nil, err
	}
	return &v1.FollowResponse{}, nil
}
func (a *authGRPCServer) UnFollow(ctx context.Context, request *v1.UnFollowRequest) (*v1.UnFollowResponse, error) {
	slog.Info("DELETE:: UnFollow")
	payloadUser, err := utils.ExtractMetadataUser(ctx)
	if err != nil {
		return nil, err
	}
	followingId, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, err
	}
	err = a.ucFollow.UnFollow(ctx, payloadUser.ID, followingId)
	if err != nil {
		return nil, err
	}
	return &v1.UnFollowResponse{}, nil
}

func (a *authGRPCServer) GetFollowers(ctx context.Context, request *v1.GetFollowersRequest) (*v1.GetFollowersResponse, error) {
	slog.Info("GET:: GetFollowers")
	followingId, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, err
	}
	followers, err := a.ucFollow.GetFollowers(ctx, followingId)
	if err != nil {
		return nil, err
	}
	return &v1.GetFollowersResponse{
		Users: lo.Map(followers, func(follower *domain.UserFollow, _ int) *v1.UserFollow {
			return &v1.UserFollow{
				Id:        follower.Id.String(),
				NickName:  follower.NickName,
				FullName:  follower.FullName,
				AvatarUrl: follower.AvatarUrl,
			}
		}),
	}, nil
}
func (a *authGRPCServer) GetFollowing(ctx context.Context, request *v1.GetFollowingRequest) (*v1.GetFollowingResponse, error) {
	slog.Info("GET:: GetFollowers")
	followerId, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, err
	}
	followings, err := a.ucFollow.GetFollowing(ctx, followerId)
	if err != nil {
		return nil, err
	}
	return &v1.GetFollowingResponse{
		Users: lo.Map(followings, func(following *domain.UserFollow, _ int) *v1.UserFollow {
			return &v1.UserFollow{
				Id:        following.Id.String(),
				NickName:  following.NickName,
				FullName:  following.FullName,
				AvatarUrl: following.AvatarUrl,
			}
		}),
	}, nil
}
