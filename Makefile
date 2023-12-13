include .test.env
export

all: build test

sqlc: 
	sqlc generate
.PHONY: sqlc

test:
	go test -v main.go

clean:
	go clean

linter-golangci: ### check by golangci linter
	golangci-lint run
.PHONY: linter-golangci

createdb:
	docker exec -it mcs-postgres createdb --username=postgres --owner=root postgres

dropdb:
	docker exec -it mcs-postgres dropdb postgres

migrateup:
	migrate -path db/migrations -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migrations -database "$(DB_URL)" -verbose down

db_docs:
	dbdocs build docs/db.dbml

db_schema:
	dbml2sql --postgres -o docs/schema.sql docs/db.dbml

wire:
	cd internal/group/app && wire && cd - &&  \
	cd internal/post/app && wire && cd - && \
	cd internal/comment/app && wire && cd - && \
	cd internal/like/app && wire && cd - && \
	cd internal/upload/app && wire && cd - && \
	cd internal/auth/app && wire && cd - && \
	cd internal/search/app && wire && cd -
.PHONY: wire

proto:
	buf generate
	
proto-gen:
	rm -f proto/gen/*.go
	rm -f third_party/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=proto/gen --go_opt=paths=source_relative \
	--go-grpc_out=proto/gen --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=proto/gen --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=third_party/swagger --openapiv2_opt=allow_merge=true,merge_file_name=go-microservices\
	proto/*.proto
	statik -src=./third_party/swagger -dest=./third_party
.PHONY: proto

docker: docker-stop docker-start
.PHONY: docker

docker-start:
	docker-compose up --build
.PHONY: docker-start

docker-stop:
	docker-compose down
.PHONY: docker-stop

docker-core: docker-core-stop docker-core-start

docker-core-start:
	docker-compose -f docker-compose-core.yaml up --build -d
.PHONY: docker-core-start

docker-core-stop:
	docker-compose -f docker-compose-core.yaml down
# --remove-orphans -v
.PHONY: docker-core-stop

docker-build:
	docker-compose down --remove-orphans -v
	docker-compose build
.PHONY: docker-build

run: run-group run-post run-comment run-like run-upload run-auth run-proxy

run-group:
	cd cmd/group && go mod tidy && go mod download && \
	CGO_ENABLED=0 go run -tags migrate github.com/dinhcanh303/go-microservices/cmd/group &
.PHONY: run-group

run-post:
	cd cmd/post && go mod tidy && go mod download && \
	CGO_ENABLED=0 go run -tags migrate github.com/dinhcanh303/go-microservices/cmd/post &
.PHONY: run-post

run-comment:
	cd cmd/comment && go mod tidy && go mod download && \
	CGO_ENABLED=0 go run -tags migrate github.com/dinhcanh303/go-microservices/cmd/comment &
.PHONY: run-comment

run-like:
	cd cmd/like && go mod tidy && go mod download && \
	CGO_ENABLED=0 go run -tags migrate github.com/dinhcanh303/go-microservices/cmd/like &
.PHONY: run-like

run-upload:
	cd cmd/upload && go mod tidy && go mod download && \
	CGO_ENABLED=0 go run -tags migrate github.com/dinhcanh303/go-microservices/cmd/upload &
.PHONY: run-upload

run-auth:
	cd cmd/auth && go mod tidy && go mod download && \
	CGO_ENABLED=0 go run -tags migrate github.com/dinhcanh303/go-microservices/cmd/auth &
.PHONY: run-auth

run-proxy:
	cd cmd/proxy && go mod tidy && go mod download && \
	CGO_ENABLED=0 go run github.com/dinhcanh303/go-microservices/cmd/proxy &
.PHONY: run-proxy

run-web:
	cd cmd/web && go mod tidy && go mod download && \
	CGO_ENABLED=0 go run github.com/dinhcanh303/go-microservices/cmd/web
.PHONY: run-web



