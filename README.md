# go-microservices

An event-driven microservices social application has been written in Golang 
## Technical stack
- Infrastructure
    - PostgreSQL
    - RabbitMQ
    - Docker and Docker-compose
    - Minio
    - Redis
## Design
![go-microservices](docs/go-microservices.svg)
## Services
No. | Service | URL
--- | ---- | -----
1 | gRPC Gateway | [http://localhost:5000](http://localhost:5000)
2 | Group Service | [http://localhost:5001](http://localhost:5001)
3 | Post Service | [http://localhost:5002](http://localhost:5002)
4 | Comment Service | [http://localhost:5003](http://localhost:5003)
5 | Like Service | [http://localhost:5004](http://localhost:5004)
6 | Upload Service | [http://localhost:5005](http://localhost:5005)<br>[http://localhost:5006](http://localhost:5006)
7 | Auth Service | [http://localhost:5007](http://localhost:5007)
8 | Notification Service | worker only
9 | Web | loading...
## Clean Architecture
![clean-architecture](docs/clean_architecture.jpg)
### Explain Clean Architecture
- Clean Architecture is a multi-layered architecture.
- Isolate Business Rules. Similar to Hexagonal and Onion architecture.
- The concentric circles represent layers, the closer to the center 
the more abstract (high level), and the further out the more detailed (low level).
- Dependency Inversion (DI) in SOLID
    High level will not depend on low level, both depend on abstraction/interface.
    Abstraction does not depend on details but vice versa.
#### Note 
- ---> In Clean Architecture arrows is Dependency Direction , not call direction.
## Clean Domain-driven Design
![clean-ddd](docs/clean_ddd.svg)

## Development
### Install tools
[Docker Desktop for Mac](https://www.docker.com/products/docker-desktop/) <br>
[TablePlus](https://tableplus.com/) or [pgAdmin4](https://www.pgadmin.org/) <br>
[Golang](https://go.dev/) <br>
[Migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) <br>
- Mac OS
```bash
brew install golang-migrate
```
- Linux
```bash
curl -L https://packagecloud.io/golang-migrate/migrate/gpgkey | apt-key add -
echo "deb https://packagecloud.io/golang-migrate/migrate/ubuntu/ $(lsb_release -sc) main" > /etc/apt/sources.list.d/migrate.list
apt-get update
apt-get install -y migrate
```
[DB Docs](https://dbdocs.io/docs)
```bash
npm install -g dbdocs
dbdocs login
```
[DBML CLI](https://dbml.dbdiagram.io/cli/#installation)
```bash
npm install -g @dbml/cli
dbml2sql --version
# or if you're using yarn
yarn global add @dbml/cli
dbml2sql --version
```
[SQLC](https://docs.sqlc.dev/en/stable/index.html)
- Mac OS
```bash
brew install sqlc
```
- Linux
```bash
sudo snap install sqlc
```
[Go mock](https://github.com/golang/mock)
```bash
go install github.com/golang/mock/mockgen@v1.6.0
```
### How to generate code
[Generate dependency injection instances with wire](https://github.com/google/wire)
```bash
make wire
```
[Generate code with sqlc](https://docs.sqlc.dev/en/stable/index.html)
```bash
make sqlc
```
[Generate proto using protobuf ](https://github.com/golang/protobuf)
```bash
make proto
```
### Documentation
Generate DB documentation
```bash
make db_docs
```
Access the DB Documentation at db [go-microservice](https://dbdocs.io/dinhcanhng303/go_microservices).Password: 123456789
### How to run
#### Run using Docker
Start docker core include (postgres , redis, rabbitmq, etc)
```bash
make docker-core
```
Start docker multi-service (group,post,like, etc)
```bash
make docker
```
[Start service traefik](https://github.com/traefik/traefik)
```bash
cd traefik
make docker
```
#### Run 
```bash
make run
```
#### Run test
```bash
make test
```
## Step By Step Create Service
- [Documents](docs/step_by_step_service.md)



