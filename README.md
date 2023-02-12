# restserver
Simple rest server with tests (unit+integration)

#IMPORTANT
- Used in-memory database for simplicity and for benifit of saving time.
- Not targetting test coverage (added samples of unit/feature tests)

Please let me know if you would like to see an implementation with any OR both of the above (as its only enhancing the current implementation, for e.g using documentDB in cloud etc).

# Notes
- Used testify to create table based tests.
- Used godog for feature tests (see folder 'features')
- Can spin up service both locally (cmd/local/main.go) OR in the cloud as lambda (cmd/lambda/main.go).

# Lambda platform considerations
- Build the lambda as linux/arm64 (aws-lambda prefered architecture)
- The maximum time a lambda can remain hot is 15 mins (more than sufficient for a REST call)
- Payload limitation is 6MB (again more than sufficient for our use-case of a tiny json payload).

# Makefile tagets
- Unit tests (make test)
- Integration tests (make test-component)
- Linting (make lint)
- debug, build.
- build-lambda

# Sample execution
```sh
// RUN SERVER
nsarup@MONSTER:~/po/restserver$ make debug
go build -tags 'production' -o build/local-server cmd/local/main.go
build/local-server
2023/02/12 14:15:41 starting service
2023/02/12 14:15:41 listening at port [:21000]

// ADD USER
nsarup@MONSTER:~$ curl -d '{"name":"user1", "email": "user1@postoffice", "password": "plaintext"}' -H "Content-Type: application/json" -X POST -w "%{http_code}\n" http://localhost:21000/users
201

// GET ALL USERS
nsarup@MONSTER:~$ curl -i -X GET http://localhost:21000/users
HTTP/1.1 200 OK
Date: Sun, 12 Feb 2023 14:20:41 GMT
Content-Length: 89
Content-Type: text/plain; charset=utf-8

[{"id":"ef3efb73-baf1-43bd-ae0d-6242a9f8644a","name":"user1","email":"user1@postoffice"}]

// GET USER BY ID
nsarup@MONSTER:~$ curl -i -X GET http://localhost:21000/users/f7a02eef-5214-4292-be89-7c1378b59df5
HTTP/1.1 200 OK
Date: Sun, 12 Feb 2023 14:23:10 GMT
Content-Length: 87
Content-Type: text/plain; charset=utf-8

{"id":"f7a02eef-5214-4292-be89-7c1378b59df5","name":"user1","email":"user1@postoffice"}
```