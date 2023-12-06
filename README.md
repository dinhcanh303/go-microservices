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
--- | --- | ---
1 | grpc-gateway | [http://localhost:5000](http://localhost:5000)
2 | group service | [http://localhost:5001](http://localhost:5001)
3 | post service | [http://localhost:5002](http://localhost:5002)
4 | comment service | [http://localhost:5003](http://localhost:5003)
5 | like service | [http://localhost:5004](http://localhost:5004)
6 | upload service | [http://localhost:5005](http://localhost:5005)<br>[http://localhost:5006](http://localhost:5006)
7 | auth service | [http://localhost:5007](http://localhost:5007)
8 | notification service | worker only
9 | web | loading...

## Clean Domain-driven Design
![clean-ddd](docs/clean_ddd.svg)

## Development

### Generate dependency injection instances with wire
```bash
> make wire
```
### Generate code with sqlc

```bash
> make sqlc
```

