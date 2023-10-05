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

wire:
	cd internal/group/app && wire
.PHONY: wire

proto:
	rm -f proto/gen/*go
	protoc --proto_path=proto --go_out=proto/gen --go_opt=paths=source_relative \
	--go-grpc_out=proto/gen --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=proto/gen --grpc-gateway_opt=paths=source_relative \
	proto/*.proto
.PHONY: proto

docker-compose: docker-compose-stop docker-compose-start
.PHONY: docker-compose

docker-compose-start:
	docker-compose up --build
.PHONY: docker-compose-start

docker-compose-stop:
	docker-compose down --remove-orphans -v
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

run: run-group run-proxy run-web

run-group:
	cd cmd/group && go mod tidy && go mod download && \
	CGO_ENABLED=0 go run -tags migrate github.com/dinhcanh303/go-microservices/cmd/group
.PHONY: run-group

run-proxy:
	cd cmd/proxy && go mod tidy && go mod download && \
	CGO_ENABLED=0 go run -tags migrate github.com/dinhcanh303/go-microservices/cmd/proxy
.PHONY: run-proxy

run-web:
	cd cmd/web && go mod tidy && go mod download && \
	CGO_ENABLED=0 go run github.com/dinhcanh303/go-microservices/cmd/web
.PHONY: run-web

