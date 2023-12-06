# Example Step By Step Create Service:
## DDD and Clean Architecture
### Folder structure
Let's focus on directories in [`/internal`](../internal/)
```
internal
├── auth
│   ├── app
│   ├── domain
│   ├── infras
│   └── usecases
├── group
│   ├── app
│   ├── domain
│   ├── infras
│   └── usecases
└── post
    ├── app
    ├── domain
    ├── infras
    └── usecases
...
```
### Example create service post 
```bash
> cd internal 
> mkdir post
> create folder domain
                infras
                usecases
                app
```
### 1.Enterprise Business Rules
#### Domain 
Contains domain models.This is the heart of the system, containing entities, states, and business logic.
##### Entity
Represents the core business objects.<br>
```go
package domain

import (
	"time"

	"github.com/dinhcanh303/go-microservices/internal/like/domain"
	sharedkernel "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID     `json:"id"`
	Status    int32         `json:"status"`
	Title     string        `json:"title"`
	Content   string        `json:"content"`
	UserID    uuid.UUID     `json:"user_id"`
	GroupID   uuid.NullUUID `json:"group_id"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}
```
##### Interfaces
Defines contracts for interactions within the domain (Domain Service Interface).<br>
```go
package domain

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/like/domain"
	sharedkernel "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
	domainUpload "github.com/dinhcanh303/go-microservices/internal/upload/domain"
	"github.com/google/uuid"
)

type (
	CommentDomainService interface {
		GetCommentsByPostID(ctx context.Context, postId uuid.UUID) ([]*sharedkernel.CommentHasChildren, error)
	}
	LikeDomainService interface {
		GetLikesByPostID(ctx context.Context, postId uuid.UUID) ([]*domain.Like, error)
	}
	UploadDomainService interface {
		GetAttachmentsByType(ctx context.Context, attachableType string, attachableId uuid.UUID) ([]*domainUpload.Attachment, error)
	}
	GroupDomainService interface {
		GetAllGroupIdByUserId(ctx context.Context, userId uuid.UUID) ([]uuid.UUID, error)
	}
	AuthDomainService interface {
		GetAllUserIdByUserId(ctx context.Context, userId uuid.UUID) ([]uuid.UUID, error)
	}
)
```
### 2.Application Business Rules
#### UseCase

![use_case](use_case.png)
##### Interfaces
```go
package posts

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/post/domain"
	"github.com/google/uuid"
)

type (
	PostRepo interface {
		Get(ctx context.Context, id uuid.UUID) (*domain.Post, error)
		Create(ctx context.Context, post *domain.Post) (*domain.Post, error)
		Update(ctx context.Context, post *domain.Post) (*domain.Post, error)
		Delete(ctx context.Context, id uuid.UUID) (bool, error)
		GetByUserId(ctx context.Context, userId uuid.UUID, limit, offset int32) ([]*domain.Post, error)
		GetByGroupId(ctx context.Context, groupId uuid.UUID, limit, offset int32) ([]*domain.Post, error)
		GetByFeed(ctx context.Context, userIds []uuid.UUID, groupIds []uuid.UUID, limit, offset int32) ([]*domain.Post, error)
	}
	UseCase interface {
		GetPost(ctx context.Context, id uuid.UUID) (*domain.Post, error)
		CreatePost(ctx context.Context, post *domain.Post) (*domain.Post, error)
		UpdatePost(ctx context.Context, post *domain.Post) (*domain.Post, error)
		DeletePost(ctx context.Context, id uuid.UUID) (bool, error)
		GetPostsByUserId(ctx context.Context, userId uuid.UUID, limit, offset int32) ([]*domain.Post, error)
		GetPostsByGroupId(ctx context.Context, groupId uuid.UUID, limit, offset int32) ([]*domain.Post, error)
		GetPostsByFeed(ctx context.Context, userIds []uuid.UUID, groupIds []uuid.UUID, limit, offset int32) ([]*domain.Post, error)
	}
)
```
###### Repository Interface
Repository Interface is defines methods for data access in the domain. 
###### UseCase Interface
UseCase Interface is defines methods representing use cases.
##### Service (UseCase)
Implements the use case interfaces.<br>
```go
package posts

import (
	"context"

	"github.com/dinhcanh303/go-microservices/internal/post/domain"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
)

type usecase struct {
	postRepo PostRepo
}

var _ UseCase = (*usecase)(nil)
var UseCaseSet = wire.NewSet(NewUseCase)

func NewUseCase(postRepo PostRepo,
) UseCase {
	return &usecase{
		postRepo: postRepo,
	}
}

// GetPostsByFeed implements UseCase.
func (uc *usecase) GetPostsByFeed(ctx context.Context, userIds []uuid.UUID, groupIds []uuid.UUID, limit int32, offset int32) ([]*domain.Post, error) {
	posts, err := uc.postRepo.GetByFeed(ctx, userIds, groupIds, limit, offset)
	if err != nil {
		return nil, errors.Wrap(err, "uc.GetPostsByFeed failed")
	}
	return posts, nil
}

// GetPostsByGroupId implements UseCase.
func (uc *usecase) GetPostsByGroupId(ctx context.Context, groupId uuid.UUID, limit int32, offset int32) ([]*domain.Post, error) {
	posts, err := uc.postRepo.GetByGroupId(ctx, groupId, limit, offset)
	if err != nil {
		return nil, errors.Wrap(err, "uc.GetPostsByGroupId failed")
	}
	return posts, nil
}
...
```
### 3.Interface Adapter
#### Infrastructure
![infras](infras.png)
##### gRPC
Implements external interface like gRPC clients.
Example Auth client.<br>
```go
package grpc

import (
	"context"
	"log/slog"

	"github.com/dinhcanh303/go-microservices/cmd/post/config"
	"github.com/dinhcanh303/go-microservices/internal/post/domain"
	"github.com/dinhcanh303/go-microservices/proto/gen"
	"github.com/google/uuid"
	"github.com/google/wire"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type authGRPCClient struct {
	conn *grpc.ClientConn
}

var _ domain.AuthDomainService = (*authGRPCClient)(nil)

var AuthGRPCClientSet = wire.NewSet(NewGRPCAuthClient)

func NewGRPCAuthClient(cfg *config.Config) (domain.AuthDomainService, error) {
	conn, err := grpc.Dial(cfg.AuthClient.URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &authGRPCClient{
		conn: conn,
	}, nil
}

// GetAllUserIdByUserId implements domain.AuthDomainService.
func (a *authGRPCClient) GetAllUserIdByUserId(ctx context.Context, userId uuid.UUID) ([]uuid.UUID, error) {
	client := gen.NewAuthServiceClient(a.conn)

	res, err := client.GetAllUserIdByUserId(ctx, &gen.GetAllUserIdByUserIdRequest{
		UserId: userId.String(),
	})
	results := make([]uuid.UUID, 0)
	if err != nil {
		slog.Warn("authGRPCClient.GetAllUserIdByUserId failed", err)
		return results, nil
	}
	for _, item := range res.UserIds {
		uuid, err := uuid.Parse(item)
		if err == nil {
			results = append(results, uuid)
		}
	}
	return results, nil
}
```
##### Postgresql
Create file query.sql
Utilizes SQLC for generating go code from SQL queries.
Example: `query.sql` for defining queries.<br>
```sql
-- name: Get :one
SELECT * FROM post.posts WHERE id = $1;

-- name: Create :one
INSERT INTO
    post.posts (
        id,
        status,
        title,
        content,
        user_id,
        group_id
    )
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: Update :one
UPDATE post.posts 
SET
    title = COALESCE(sqlc.narg(title),title),
    content = COALESCE(sqlc.narg(content),content),
    status = COALESCE(sqlc.narg(status),status)
WHERE id = sqlc.arg(id) RETURNING *;

-- name: GetByGroupId :many
SELECT * FROM post.posts WHERE group_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: GetByUserId :many
SELECT * FROM post.posts 
WHERE user_id = $1 AND group_id IS NULL ORDER BY created_at DESC LIMIT $2 OFFSET $3;

-- name: GetByFeed :many
SELECT *
FROM post.posts
WHERE user_id = ANY($1::uuid[])
   OR group_id = ANY($2::uuid[])
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;
-- name: Delete :exec
DELETE FROM post.posts WHERE id = $1;
```
SQLC generation using the `make sqlc` command:
```bash
> make sqlc 
```
Generated files:db.go,models.go and query.sql.go.<br>
![sqlc_gen](sqlc_gen.png)
##### Repository
Implements repository interface for data access.<br>
```go
package repo

import (
	"context"
	"database/sql"

	"github.com/dinhcanh303/go-microservices/internal/post/domain"
	"github.com/dinhcanh303/go-microservices/internal/post/infras/postgresql"
	"github.com/dinhcanh303/go-microservices/internal/post/usecases/posts"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type postRepo struct {
	pg postgres.DBEngine
}

func NewPostRepo(pg postgres.DBEngine) posts.PostRepo {
	return &postRepo{pg: pg}
}

var _ posts.PostRepo = (*postRepo)(nil)

var RepositoryPostSet = wire.NewSet(NewPostRepo)

// GetByFeed implements posts.PostRepo.
func (rp *postRepo) GetByFeed(ctx context.Context, userIds []uuid.UUID, groupIds []uuid.UUID, limit int32, offset int32) ([]*domain.Post, error) {
	db := rp.pg.GetDB()
	querier := postgresql.New(db)

	results, err := querier.GetByFeed(ctx, postgresql.GetByFeedParams{
		Column1: userIds,
		Column2: groupIds,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, errors.Wrap(err, "qtx.GetByFeed(ctx, userIds, groupIds, limit, offset) failed")
	}
	return lo.Map(results, func(item postgresql.PostPost, _ int) *domain.Post {
		return &domain.Post{
			ID:        item.ID,
			Title:     item.Title,
			Content:   item.Content,
			Status:    item.Status,
			UserID:    item.UserID,
			GroupID:   item.GroupID,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}
	}), nil
}
...
```
### 4.Framework and Drivers
#### Application
##### Router
Using gRPC server (router) for handling inbound call UseCase (service) and Domain Service Interface.<br> 
```go
package router

import (
	"context"

	"github.com/dinhcanh303/go-microservices/cmd/post/config"
	domainComment "github.com/dinhcanh303/go-microservices/internal/comment/domain"
	domainLike "github.com/dinhcanh303/go-microservices/internal/like/domain"
	sharedkernel "github.com/dinhcanh303/go-microservices/internal/pkg/shared_kernel"
	"github.com/dinhcanh303/go-microservices/internal/post/domain"
	"github.com/dinhcanh303/go-microservices/internal/post/usecases/posts"
	domainUpload "github.com/dinhcanh303/go-microservices/internal/upload/domain"
	"github.com/dinhcanh303/go-microservices/pkg/utils"
	"github.com/dinhcanh303/go-microservices/proto/gen"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type postGRPCServer struct {
	gen.UnimplementedPostServiceServer
	cfg                  *config.Config
	uc                   posts.UseCase
	uploadDomainService  domain.UploadDomainService
	commentDomainService domain.CommentDomainService
	likeDomainService    domain.LikeDomainService
	groupDomainService   domain.GroupDomainService
	authDomainService    domain.AuthDomainService
}

var _ gen.PostServiceServer = (*postGRPCServer)(nil)

var PostGRPCServerSet = wire.NewSet(NewGRPCPostServer)

func NewGRPCPostServer(
	grpcServer *grpc.Server,
	cfg *config.Config,
	uc posts.UseCase,
	uploadDomainService domain.UploadDomainService,
	commentDomainService domain.CommentDomainService,
	likeDomainService domain.LikeDomainService,
	groupDomainService domain.GroupDomainService,
	authDomainService domain.AuthDomainService,
) gen.PostServiceServer {
	svc := postGRPCServer{
		cfg:                  cfg,
		uc:                   uc,
		uploadDomainService:  uploadDomainService,
		commentDomainService: commentDomainService,
		likeDomainService:    likeDomainService,
		groupDomainService:   groupDomainService,
		authDomainService:    authDomainService,
	}
	gen.RegisterPostServiceServer(grpcServer, &svc)
	reflection.Register(grpcServer)
	return &svc
}

func (p *postGRPCServer) CreatePost(ctx context.Context, request *gen.CreatePostRequest) (*gen.CreatePostResponse, error) {
	slog.Info("POST: CreatePost")
	user, err := utils.ExtractMetadataUser(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Extract Metadata User failed")
	}
	groupId, _ := uuid.Parse(request.Post.GroupId)
	model := domain.Post{
		Title:   request.Post.Title,
		Content: request.Post.Content,
		Status:  request.Post.Status,
		UserID:  user.ID,
		GroupID: uuid.NullUUID{
			UUID:  groupId,
			Valid: request.Post.GroupId != "",
		},
	}
	post, err := p.uc.CreatePost(ctx, &model)
	if err != nil {
		return nil, errors.Wrap(err, "uc.CreatePost failed")
	}
	res := &gen.CreatePostResponse{
		Post: &gen.PostResponse{
			Id:        post.ID.String(),
			Title:     post.Title,
			Content:   post.Content,
			UserId:    post.UserID.String(),
			GroupId:   post.GroupID.UUID.String(),
			Status:    post.Status,
			CreatedAt: timestamppb.New(post.CreatedAt),
			UpdatedAt: timestamppb.New(post.UpdatedAt),
		},
	}
	return res, nil
}
...
```
##### App
Contains the main application logic and wiring.<br>
```go
package app

import (
	"github.com/dinhcanh303/go-microservices/cmd/post/config"
	"github.com/dinhcanh303/go-microservices/internal/post/domain"
	"github.com/dinhcanh303/go-microservices/internal/post/usecases/posts"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/dinhcanh303/go-microservices/proto/gen"
)

type App struct {
	Cfg              *config.Config
	PG               postgres.DBEngine
	UC               posts.UseCase
	PostGRPCServer   gen.PostServiceServer
	CommentDomainSvc domain.CommentDomainService
	LikeDomainSvc    domain.LikeDomainService
	UploadDomainSvc  domain.UploadDomainService
}

func New(
	cfg *config.Config,
	pg postgres.DBEngine,
	uc posts.UseCase,
	postGRPCServer gen.PostServiceServer,
	commentDomainSvc domain.CommentDomainService,
	likeDomainSvc domain.LikeDomainService,
	uploadDomainSvc domain.UploadDomainService) *App {
	return &App{
		Cfg:              cfg,
		UC:               uc,
		PG:               pg,
		PostGRPCServer:   postGRPCServer,
		CommentDomainSvc: commentDomainSvc,
		LikeDomainSvc:    likeDomainSvc,
		UploadDomainSvc:  uploadDomainSvc,
	}
}
```
##### Wire (Dependency Injection (DI) of Google)
Utilizes Google Wire for dependency injection.<br>
```go
//go:build wireinject
// +build wireinject

package app

import (
	"github.com/dinhcanh303/go-microservices/cmd/post/config"
	"github.com/dinhcanh303/go-microservices/internal/post/app/router"
	infrasGRPC "github.com/dinhcanh303/go-microservices/internal/post/infras/grpc"
	"github.com/dinhcanh303/go-microservices/internal/post/infras/repo"
	postsUC "github.com/dinhcanh303/go-microservices/internal/post/usecases/posts"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/google/wire"
	"google.golang.org/grpc"
)

func InitApp(
	cfg *config.Config,
	dbConnStr postgres.DBConnString,
	grpcServer *grpc.Server,
) (*App, func(), error) {
	panic(wire.Build(
		New,
		dbEngineFunc,
		router.PostGRPCServerSet,
		repo.RepositoryPostSet,
		postsUC.UseCaseSet,
		infrasGRPC.CommentGRPCClientSet,
		infrasGRPC.LikeGRPCClientSet,
		infrasGRPC.UploadGRPCClientSet,
		infrasGRPC.GroupGRPCClientSet,
		infrasGRPC.AuthGRPCClientSet,
	))
}
func dbEngineFunc(url postgres.DBConnString) (postgres.DBEngine, func(), error) {
	db, err := postgres.NewPostgresDB(url)
	if err != nil {
		return nil, nil, err
	}
	return db, func() { db.Close() }, nil
}
```
Then using wire generation file `wire_gen.go` (file code gen dependency injection).
```bash
> make wire
```
#### Cmd
##### Folder structure
Let's focus on directories in [`/cmd`](../cmd/)
```
internal
├── auth
│   ├── config
│   ├── config.yml
│   ├── main.go
├── group
│   ├── config
│   ├── config.yml
│   ├── main.go
└── post
    ├── config
    ├── config.yml
    ├── main.go
...
```
##### Example Post
###### Config
- `config.go` using mapping config.yml and environment from docker or kubectl configuration
```go
package config

import (
	"fmt"
	"log"
	"os"

	configs "github.com/dinhcanh303/go-microservices/pkg/config"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		configs.App   `yaml:"app"`
		configs.HTTP  `yaml:"http"`
		configs.Log   `yaml:"logger"`
		PG            `yaml:"postgres"`
		CommentClient `yaml:"comment_client"`
		LikeClient    `yaml:"like_client"`
		UploadClient  `yaml:"upload_client"`
		GroupClient   `yaml:"group_client"`
		AuthClient    `yaml:"auth_client"`
	}

	PG struct {
		PoolMax int    `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
		DsnURL  string `env-required:"true" yaml:"dsn_url" env:"PG_DSN_URL"`
	}

	CommentClient struct {
		URL string `env-required:"true" yaml:"comment_url" env:"COMMENT_CLIENT_URL"`
	}
	LikeClient struct {
		URL string `env-required:"true" yaml:"like_url" env:"LIKE_CLIENT_URL"`
	}
	UploadClient struct {
		URL string `env-required:"true" yaml:"upload_url" env:"UPLOAD_CLIENT_URL"`
	}
	GroupClient struct {
		URL string `env-required:"true" yaml:"upload_url" env:"GROUP_CLIENT_URL"`
	}
	AuthClient struct {
		URL string `env-required:"true" yaml:"upload_url" env:"AUTH_CLIENT_URL"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// debug
	fmt.Println("config path: " + dir)

	err = cleanenv.ReadConfig(dir+"/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

```
- `config.yml`
Config information app , http , database ,logger and call service though gRPC ,REST ,Message Queues,...
```yml
app:
  name: 'post-service'
  version: '1.0.0'

http:
  host: '0.0.0.0'
  port: 5002

postgres:
  pool_max: 2
  dsn_url: host=127.0.0.1 user=postgres password=P@ssw0rd dbname=postgres sslmode=disable

group_client:
  group_url: 0.0.0.0:5001
comment_client:
  comment_url: 0.0.0.0:5003
like_client:
  like_url: 0.0.0.0:5004
upload_client:
  upload_url: 0.0.0.0:5006
auth_client:
  auth_url: 0.0.0.0:5007

logger:
  log_level: 'debug'
  rollbar_env: 'post-service'
```
##### main.go (server)
File run server gRPC ,Echo, Gin,... config call inbound application
```go
package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/dinhcanh303/go-microservices/cmd/post/config"
	"github.com/dinhcanh303/go-microservices/internal/post/app"
	"github.com/dinhcanh303/go-microservices/pkg/logger"
	"github.com/dinhcanh303/go-microservices/pkg/postgres"
	"github.com/sirupsen/logrus"
	"go.uber.org/automaxprocs/maxprocs"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
)

func main() {
	_, err := maxprocs.Set()
	if err != nil {
		slog.Error("Failed set max process", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error("Failed get config", err)
	}
	slog.Info("⚡ Init App", "name", cfg.Name, "version", cfg.Version)

	//set up logrus
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logger.ConvertLogLevel(cfg.Log.Level))

	//integrate Logrus with the slog logger
	logrusHandle := logger.NewLogrusHandler(logrus.StandardLogger())
	slog.New(logrusHandle)

	server := grpc.NewServer()

	go func() {
		defer server.GracefulStop()
		<-ctx.Done()
	}()
	cleanup := prepareApp(ctx, cancel, cfg, server)

	//gRPC Server
	address := fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port)
	network := "tcp"
	l, err := net.Listen(network, address)
	if err != nil {
		slog.Error("Failed to listen to address", err, "Network", network, "Address", address)
		cancel()
		<-ctx.Done()
	}
	slog.Info("🌏 start server...", "address", address)
	defer func() {
		if err1 := l.Close(); err != nil {
			slog.Error("failed to close", err1, "network", network, "address", address)
			<-ctx.Done()
		}
	}()
	err = server.Serve(l)
	if err != nil {
		slog.Error("failed start gRPC server", err, "network", network, "address", address)
		cancel()
		<-ctx.Done()
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	select {
	case v := <-quit:
		cleanup()
		slog.Info("signal.Notify", v)
	case done := <-ctx.Done():
		cleanup()
		slog.Info("ctx.Done", "app done", done)
	}

}

func prepareApp(ctx context.Context, cancel context.CancelFunc, cfg *config.Config, server *grpc.Server) func() {
	_, cleanup, err := app.InitApp(cfg, postgres.DBConnString(cfg.PG.DsnURL), server)
	if err != nil {
		slog.Error("Failed init app", err)
		cancel()
		<-ctx.Done()
	}
	return cleanup
}

```
### Generation protocol
#### Folder structure
Let's focus on directories in [`/proto`](../proto/)
```
internal
├── gen
├── google
└── pb
└── protoc-gen-openapiv2
└── validate
└── auth.proto
└── buf.yaml
└── commom.proto
...
```
[Doc:grpc-ecosystem/grpc-gateway/v2](https://github.com/grpc-ecosystem/grpc-gateway)
The gRPC-Gateway is a plugin of the Google protocol buffers compiler protoc. It reads protobuf service definitions and generates a reverse-proxy server which translates a RESTful HTTP API into gRPC. This server is generated according to the google.api.http annotations in your service definitions.

This helps you provide your APIs in both gRPC and RESTful style at the same time.
![gRPC-Gateway](architecture_introduction_diagram.svg)
##### 1.Create file `post.proto` using syntax proto 3
```proto
syntax="proto3";

package post;

import "google/api/annotations.proto";
import "protoc-gen-openapiv2/options/annotations.proto";
import "google/protobuf/timestamp.proto";
import "upload.proto";
import "like.proto";
import "comment.proto";

option go_package = "github.com/dinhcanh303/go-microsevices/proto/gen";

service PostService {
    rpc CreatePost(CreatePostRequest) returns (CreatePostResponse){
        option (google.api.http) = {
            post: "/api/v1/posts"
            body: "*"
        };
        option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {
            summary: "Create Post"
            description: ""
        };
    }
}
message CreatePostRequest {
    PostRequest post = 1;
}
message CreatePostResponse {
    PostResponse post = 1;
}
message PostRequest {
    string title = 1;
    string content = 2;
    int32 status = 3;
    string group_id = 4;
}
message PostResponse {
    string id = 1;
    string title = 2;
    string content = 3;
    int32 status = 4;
    string user_id = 5;
    string group_id = 6;
    google.protobuf.Timestamp created_at = 7;
    google.protobuf.Timestamp updated_at = 8;
}
```
##### 2. buf.yaml config 
```yaml
version: v1
name: buf.build/dinhcanh303/go-microservices
deps:
  - buf.build/googleapis/googleapis
  - buf.build/grpc-ecosystem/grpc-gateway
lint:
  use:
    - DEFAULT
  ignore_only:
    PACKAGE_DIRECTORY_MATCH:
      - common.proto
      - post.proto <--
    PACKAGE_VERSION_SUFFIX:
      - common.proto
      - post.proto <--
    RPC_REQUEST_RESPONSE_UNIQUE:
      - common.proto
      - post.proto <--
    RPC_RESPONSE_STANDARD_NAME:
      - common.proto
      - post.proto <--
```
Then using `make proto` let the tool automatically generate the necessary files
```bash
> make proto
```



