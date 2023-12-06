# Step By Step create service:


## DDD and Clean Architecture
### Example create service post 
```bash
> cd internal 
> mkdir post
> create folder domain
                infras
                usecases
                app
```
![folder_post](folder_post.png)
### 1.Enterprise Business Rules
#### Domain 
- Contains domain models.This is the heart of the system, containing entities, states, and business logic.
##### Entity
![domain_entity](domain_entity.png)
##### Interfaces
![domain_interface](domain_interface.png)
### 2.Application Business Rules
#### UseCase
![use_case](use_case.png)
##### Interfaces
![use_case_interface](use_case_interface.png)
###### Repository Interface
- Repository Interface 
###### UseCase Interface
- UseCase Interface 
##### Service (UseCase)
![use_case_service](use_case_service.png)

### 3.Interface Adapter
#### Infrastructure
![infras](infras.png)
##### gRPC
- Example Auth client 
![grpc_auth_client](grpc_auth_client.png)
...
##### Postgresql
- Create file query.sql written query sql using syntax sqlc support
![postgres_query](postgres_query.png)
- Then using command:
```bash
> make sqlc 
```
![sqlc_gen](sqlc_gen.png)
- sqlc generation multi file (db.go,models.go and query.sql.go)
##### Repo
![repo](repo.png)

