include .env
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

env:
	docker run --env-file .env upload

wire:
	cd internal/group/app && wire && cd - && \
	cd internal/post/app && wire && cd - && \
	cd internal/comment/app && wire && cd - && \
	cd internal/like/app && wire && cd - && \
	cd internal/upload/app && wire && cd -
.PHONY: wire

proto-gen:
	buf generate
	
proto:
	rm -f proto/gen/*.go
	rm -f third_party/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=proto/gen --go_opt=paths=source_relative \
	--go-grpc_out=proto/gen --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=proto/gen --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=third_party/swagger --openapiv2_opt=allow_merge=true,merge_file_name=go-microservices\
	proto/*.proto
	statik -src=./third_party/swagger -dest=./third_party
.PHONY: proto

docker-compose: docker-compose-stop docker-compose-start
.PHONY: docker-compose

docker-compose-start:
	docker-compose up --build
.PHONY: docker-compose-start

docker-compose-stop:
	docker-compose down
.PHONY: docker-compose-stop

docker-compose-core: docker-compose-core-stop docker-compose-core-start

docker-compose-core-start:
	docker-compose -f docker-compose-core.yaml up --build
.PHONY: docker-compose-core-start

docker-compose-core-stop:
	docker-compose -f docker-compose-core.yaml down --remove-orphans -v
.PHONY: docker-compose-core-stop

docker-compose-build:
	docker-compose down --remove-orphans -v
	docker-compose build
.PHONY: docker-compose-build

run: run-group run-post run-proxy run-web

run-group:
	cd cmd/group && go mod tidy && go mod download && \
	CGO_ENABLED=0 go run -tags migrate github.com/dinhcanh303/go-microservices/cmd/group
.PHONY: run-group

run-post:
	cd cmd/post && go mod tidy && go mod download && \
	CGO_ENABLED=0 go run -tags migrate github.com/dinhcanh303/go-microservices/cmd/post
.PHONY: run-post

run-proxy:
	cd cmd/proxy && go mod tidy && go mod download && \
	CGO_ENABLED=0 go run -tags migrate github.com/dinhcanh303/go-microservices/cmd/proxy
.PHONY: run-proxy

run-web:
	cd cmd/web && go mod tidy && go mod download && \
	CGO_ENABLED=0 go run github.com/dinhcanh303/go-microservices/cmd/web
.PHONY: run-web



