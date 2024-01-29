package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/dinhcanh303/go-microservices/internal/auth/domain"
	"github.com/dinhcanh303/go-microservices/internal/auth/usecases/keys"
	"github.com/dinhcanh303/go-microservices/internal/pkg/event"
	"github.com/dinhcanh303/go-microservices/pkg/constant"
	"github.com/dinhcanh303/go-microservices/pkg/ldap"
	"github.com/dinhcanh303/go-microservices/pkg/redis"
	"github.com/dinhcanh303/go-microservices/pkg/token"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type service struct {
	repo                UserRepo
	ucKeys              keys.UseCase
	ldapClient          ldap.LdapClient
	jwt                 token.JWT
	userCreatedEventPub UserCreatedEventPublisher
	userDeletedEventPub UserDeletedEventPublisher
	redis               redis.RedisEngine
}

var _ UseCase = (*service)(nil)

func NewUseCase(
	repo UserRepo,
	ucKeys keys.UseCase,
	ldapClient ldap.LdapClient,
	jwt token.JWT,
	userCreatedEventPub UserCreatedEventPublisher,
	userDeletedEventPub UserDeletedEventPublisher,
	redis redis.RedisEngine,
) UseCase {
	return &service{
		repo:                repo,
		ucKeys:              ucKeys,
		ldapClient:          ldapClient,
		jwt:                 jwt,
		userCreatedEventPub: userCreatedEventPub,
		userDeletedEventPub: userDeletedEventPub,
		redis:               redis,
	}
}

var UseCaseSet = wire.NewSet(NewUseCase)

var CACHE_SV_AUTH_USERS = "sv_auth_users"

// GetUsersBirthDayByCurrentDay implements UseCase.
func (s *service) GetUsersBirthDayByCurrentDay(ctx context.Context) ([]*domain.User, error) {
	users, err := s.repo.GetUsersBirthDayByCurrentDay(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetUsersBirthDayByCurrentMonth implements UseCase.
func (s *service) GetUsersBirthDayByCurrentMonth(ctx context.Context) ([]*domain.User, error) {
	users, err := s.repo.GetUsersBirthDayByCurrentMonth(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetUsersInviteGroup implements UseCase.
func (s *service) GetUsersInviteGroup(ctx context.Context, groupIds []uuid.UUID, limit int32, offset int32) ([]*domain.User, error) {
	users, err := s.repo.GetUsersInviteGroup(ctx, groupIds, limit, offset)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// UpdateUser implements UseCase.
func (s *service) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	user, err := s.repo.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	//Del cache
	err = s.redis.Invalidate(CACHE_SV_AUTH_USERS)
	if err != nil {
		slog.Error("Invalidate cache list users failed")
	}
	return user, nil
}

// GetUsers implements UseCase.
func (s *service) GetUsers(ctx context.Context, search string, limit int32, offset int32) ([]*domain.User, error) {
	if limit == 0 {
		limit = 1000
	}
	var users []*domain.User
	err := utils.HandleHitCache(users, s.redis, CACHE_SV_AUTH_USERS)
	if err != nil {
		users, err = s.repo.GetUsers(ctx, search, limit, offset)
		if err != nil {
			return nil, errors.Wrap(err, "uc.GetUsers failed")
		}
		err = s.redis.Set(CACHE_SV_AUTH_USERS, users, 0)
		if err != nil {
			return nil, errors.Wrap(err, "failed set value in cache")
		}
	}
	return users, nil
}

// GetUser implements UseCase.
func (s *service) GetUser(ctx context.Context, userId uuid.UUID) (*domain.User, error) {
	user, err := s.repo.GetUser(ctx, userId)
	//
	if err != nil {
		return nil, err
	}
	return user, nil
}

// HandleRefreshToken implements UseCase.
func (s *service) HandleRefreshToken(ctx context.Context, email, refreshToken string) (*domain.UserAuth, error) {
	foundUser, err := findUserByEmail(s.repo, ctx, email)
	if err != nil {
		return nil, err
	}
	results, err := createTokenPairAndResponse(ctx, s, foundUser, false)
	if err != nil {
		return nil, err
	}
	s.ucKeys.CreateKeyToken(ctx, &domain.Key{
		UserID:       foundUser.ID,
		RefreshToken: results.RefreshToken,
		// RefreshTokensUsed: refreshToken,
	})
	return nil, nil
}

// GetAllUserIdByUserId implements UseCase.
func (s *service) GetUserIdsOfCompanyByUserId(ctx context.Context, userId uuid.UUID) ([]uuid.UUID, error) {
	user, err := s.repo.GetUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	suffixEmailCompany := constant.SuffixEmailCompany
	if strings.Contains(user.Email, suffixEmailCompany) {
		userIds, err := s.repo.GetUserIdsOfCompany(ctx, suffixEmailCompany)
		if err != nil {
			return nil, err
		}
		return userIds, nil
	}
	return nil, nil
}

// Logout implements UseCase.
func (s *service) Logout(ctx context.Context, key *domain.Key) error {
	err := s.ucKeys.DeleteKeyByID(ctx, key.ID)
	return err
}

// SignIn implements UseCase.
func (s *service) SignIn(ctx context.Context, email string, password string) (*domain.UserAuth, error) {
	slog.Info("Service Auth:: SignIn")
	isEmailCompany := checkEmailCompany(email)
	if isEmailCompany {
		auth, _, err := s.ldapClient.Authenticate(email, password)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		if auth {
			foundUser, _ := findUserByEmail(s.repo, ctx, email)
			if foundUser == nil {
				hashPassword, err := utils.HashPassword(password)
				if err != nil {
					return nil, errors.New("create account using ldap failed because hash password failed")
				}
				name := strings.Split(email, "@")[0]
				foundUser, err = s.repo.CreateUser(ctx, &domain.User{
					ID:        uuid.New(),
					Email:     email,
					Password:  hashPassword,
					FirstName: name,
					LastName:  name,
					FullName:  name,
					NickName:  name,
				})
				if err != nil {
					return nil, errors.Wrap(err, "create account using ldap failed")
				}
				//Del cache
				err = s.redis.Invalidate(CACHE_SV_AUTH_USERS)
				if err != nil {
					slog.Error("Invalidate cache list users failed")
				}
			}
			return createTokenPairAndResponse(ctx, s, foundUser, false)
		}
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	foundUser, err := findUserByEmail(s.repo, ctx, email)
	if err != nil {
		return nil, err
	}
	match, err := utils.ComparePassword(foundUser.Password, password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	if !match {
		return nil, status.Error(codes.Unauthenticated, "Password mismatch")
	}
	return createTokenPairAndResponse(ctx, s, foundUser, false)
}

// SignUp implements UseCase.
func (s *service) SignUp(ctx context.Context, email, password, fistName, lastName string) (*domain.UserAuth, error) {
	slog.Info("Service Auth:: SignUp")
	isEmailCompany := checkEmailCompany(email)
	if isEmailCompany {
		return nil, status.Error(codes.AlreadyExists, fmt.Sprintf("Email has suffix is %s does't created. Please use the login function for this type of email!", constant.SuffixEmailCompany))
	}
	foundUser, _ := s.repo.GetUserByEmail(ctx, email)
	if foundUser != nil {
		return nil, status.Error(codes.AlreadyExists, "User already exists")
	}
	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}
	model := domain.NewUser(email, passwordHash, fistName, lastName)
	newUser, err := s.repo.CreateUser(ctx, model)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	//Del cache
	err = s.redis.Invalidate(CACHE_SV_AUTH_USERS)
	if err != nil {
		slog.Error("Invalidate cache list users failed")
	}
	if err == nil {
		// Publish event created group
		eventBytes, err := json.Marshal(event.UserCreated{
			ID:     newUser.ID,
			Name:   newUser.FullName,
			Avatar: "",
			Email:  newUser.Email,
			Type:   "user",
		})
		if err != nil {
			slog.Error("json marshal error", err)
		}
		s.userCreatedEventPub.Publish(ctx, eventBytes, "text/plain")
	}
	return createTokenPairAndResponse(ctx, s, newUser, true)
}

func checkEmailCompany(email string) bool {
	return strings.Contains(email, constant.SuffixEmailCompany)
}
func findUserByEmail(repo UserRepo, ctx context.Context, email string) (*domain.User, error) {
	foundUser, err := repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, status.Error(codes.NotFound, "User not registered")
	}
	return foundUser, err
}
func createKeyPair() (string, string, error) {
	publicKey, err := utils.GenerateRandomHexBytes(64)
	if err != nil {
		return "", "", err
	}
	privateKey, err := utils.GenerateRandomHexBytes(64)
	if err != nil {
		return "", "", err
	}
	return publicKey, privateKey, nil
}

func createTokenPair(jwt token.JWT, user *domain.User, publicKey, privateKey string) (string, string, error) {
	payload := token.NewPayload(user.ID,
		user.Email,
		user.FullName,
		user.Role, "",
		2*24*time.Hour)
	accessToken, err := jwt.CreateToken(payload, publicKey)
	if err != nil {
		return "", "", err
	}
	payload = &token.Payload{
		ExpiredAt: payload.ExpiredAt.Add(5 * 24 * time.Hour),
	}
	refreshToken, err := jwt.CreateToken(payload, privateKey)
	if err != nil {
		return "", "", err
	}
	_, err = jwt.VerifyToken(accessToken, publicKey)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}
func createTokenPairAndResponse(ctx context.Context, service *service, user *domain.User, isSingUp bool) (*domain.UserAuth, error) {
	publicKey, privateKey, err := createKeyPair()
	if err != nil {
		return nil, status.Error(codes.Unknown, "Create Public Private Key failed")
	}
	accessToken, refreshToken, err := createTokenPair(service.jwt, user, publicKey, privateKey)
	if err != nil {
		return nil, status.Error(codes.Unknown, "Create Token Failed")
	}
	model := &domain.Key{
		UserID:     user.ID,
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}
	if !isSingUp {
		model.RefreshToken = refreshToken
	}
	_, err = service.ucKeys.CreateKeyToken(ctx, model)
	if err != nil {
		return nil, status.Error(codes.Unknown, "Create Key Store Failed")
	}
	return &domain.UserAuth{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
