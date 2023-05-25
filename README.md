## Merkle Tree API
This is Go implementation of User access microservice,

We Postgres for storage and Kafka for Asynchronous notification

## GRPC API
This service provide a GRPC API at port :50051
````go

service UserService {
rpc AddUser(UserRequest) returns (UserResponse) {}
rpc UpdateUser(UserRequest) returns (UserResponse) {}
rpc RemoveUser(UserRequest) returns (UserResponse) {}
rpc GetUsers(GetUsersRequest) returns (GetUsersResponse) {}
}

message UserRequest {
string id = 1;
string first_name = 2;
string last_name = 3;
string nickname = 4;
string password = 5;
string email = 6;
string country = 7;
}
````

## Rest Endpoints
The API has below endpoints:

### POST /users
Create a new User
This endpoint accepts a JSON payload containing a user data.

Request Payload

Example request payload:
````json
{
  "first_name": "Mohamed",
  "last_name": "Aly",
  "nickname": "MA",
  "password": "pw2",
  "email": "ma@example.com",
  "country": "FIN"
}
````

### GET /users
Return a paginated list of Users, allowing for filtering by certain criteria (e.g. all Users with the country "UK")


### PUT /users/:id
Update a user
Request Payload

````json
{
  "first_name": "Mohamed",
  "last_name": "Aly",
  "nickname": "MA",
  "password": "pw2",
  "email": "ma@example.com",
  "country": "FIN"
}
````

### DELETE /users/:id
Delete a user

### GET /health
health check


### Project layout

This layout is following pattern:

```text
merkleTree
└───
    ├── .github
    │   └── workflows
    │     └── go.yml
    ├── cmd
    │   └── main.go
    │   └── app
    │     └── setup.go
    │     └── app.go
    │     └── context.go
    ├── internal
    │   └── app
    │     └── handler.go
    │     └── user_service.go
    │   └── notifier
    │     └── kafka_notifier.go
    │   └── repository
    │     └── user_repository.go
    │   └── domain
    │     └── user.go
    ├── pkg
    │   └── grpc
    │     └── grpcServer.go
    │     └── user.pb.go
    │     └── user._grpc.pb.go
    │     └── user.proto
    ├── build
    │   └── Dockerfile
    ├── Makefile
    ├── README.md
    └── <source packages>
```

## Setup

### Getting started
user-svc is available in github
[user-svc](https://github.com/reactivejson/users-svc)

```shell
go get github.com/reactivejson/users-svc
```

#### Run
```shell
make environment-start
go run cmd/main.go
```

#### Build
```shell
make build
```
#### Testing
```shell
make test
```

#### Integration Tests
```shell
make environment-start
make test-integration
```
### Build docker image:

```bash
make docker-build
```
This will build this application docker image so-called user-svc
