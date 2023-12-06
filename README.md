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

## Clean Domain-driven Design
![clean-ddd](docs/clean_ddd.svg)

## Development

### Generate dependency injection instances with wire
- Compile-time Dependency Injection for Go of Google
- [wire](https://github.com/google/wire)
```bash
> make wire
```
### Generate code with sqlc
- [sqlc](https://docs.sqlc.dev/en/stable/index.html)
```bash
> make sqlc
```
### Generate proto using protobuf 
- [protobuf](https://github.com/golang/protobuf)
```bash
> make proto-gen
```
## Start go-microservices
### Start docker core include (postgres , redis, rabbitmq, etc)
```bash
> make docker-compose-core
```
### Start docker multi-service (group,post,like, etc)
```bash
> make docker-compose
```
### Start service traefik 
- [traefik](https://github.com/traefik/traefik)
```bash
> cd traefik -> make docker-compose
```



