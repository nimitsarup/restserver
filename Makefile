BINPATH ?= build

.PHONY: lint
lint:
	golangci-lint run ./... --timeout 2m --tests=false --skip-dirs=features

.PHONY: test-component
test-component:
	go test -cover -coverpkg=github.com/nimitsarup/restserver/... -component

.PHONY: test
test:
	go test -race -cover ./...

.PHONY: build
build:
	go build -tags 'production' -o $(BINPATH)/local-server cmd/local/main.go

.PHONY: debug
debug:
	go build -tags 'production' -o $(BINPATH)/local-server cmd/local/main.go
	$(BINPATH)/local-server

.PHONY: build-lambda-archive
build-lambda-archive:
	env GOOS=linux GOARCH=amd64 go build -o $(BINPATH)/lambda-server cmd/lambda/main.go
	zip -j $(BINPATH)/lambda-server.zip $(BINPATH)/lambda-server

