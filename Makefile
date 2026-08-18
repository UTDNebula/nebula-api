
EXEC_NAME?=rest-api
DOCKER_RUNNER?=docker
REGISTRY?=localhost:5000
RELEASE_TAG=$(shell git rev-parse --short HEAD 2>/dev/null || echo "latest")
BASE_TAG=$(REGISTRY)/utdnebula/rest/rest-api

.PHONY: all setup check test docs docs-rest build build-rest clean docker docker-rest

all: check test build

setup:
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/swaggo/swag/cmd/swag@latest

check:
	go mod tidy
	go vet ./...
	staticcheck ./...
	gofmt -w .
	goimports -w .

test:
	go test ./... -count=1

docs-rest:
	swag fmt -d rest
	swag init -d rest -g server.go -o rest/docs --outputTypes yaml,go

docs: docs-rest

build-rest: docs-rest
	go build -o $(EXEC_NAME) ./rest

build: build-rest

clean:
	rm -f $(EXEC_NAME) rest/$(EXEC_NAME) $(EXEC_NAME).exe rest/$(EXEC_NAME).exe

docker-rest:
	$(DOCKER_RUNNER) build -f rest/Dockerfile -t $(BASE_TAG):$(RELEASE_TAG) .
	$(DOCKER_RUNNER) tag $(BASE_TAG):$(RELEASE_TAG) $(BASE_TAG):latest

docker: docker-rest

