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

run: run-web run-proxy

run-group:
	cd cmd/group && go mod tidy && go mod dowload && \
	CGO_ENABLED=0 go run github.com/dinhcanh303/go-microservices/cmd/group
.PHONY: run-group

run-proxy:
	cd cmd/proxy && go mod tidy && go mod download && \
	CGO_ENABLED=0 go run -tags migrate github.com/dinhcanh303/go-coffeeshop/cmd/proxy
.PHONY: run-proxy

run-web:
	cd cmd/web && go mod tidy && go mod download && \
	CGO_ENABLED=0 go run github.com/dinhcanh303/go-coffeeshop/cmd/web
.PHONY: run-web

