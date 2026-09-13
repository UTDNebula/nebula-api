
RELEASE_TAG=$(shell git rev-parse --short HEAD 2>/dev/null || echo "latest")
DOCKER_RUNNER?=docker

REST_EXEC_NAME?=rest-api
REST_REGISTRY?=localhost:5000
REST_BASE_TAG=$(REST_REGISTRY)/utdnebula/rest/rest-api

GRAPH_EXEC_NAME?=graph-api
GRAPH_REGISTRY?=localhost:5001
GRAPH_BASE_TAG=$(GRAPH_REGISTRY)/utdnebula/graphql/graphql-api

all: check test build

setup:
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/swaggo/swag/cmd/swag@latest

format:
	go mod tidy
	gofmt -w .
	goimports -w .

check:
	go mod tidy
	git diff --exit-code -- go.mod go.sum
	go vet ./...
	staticcheck ./...
	@test -z "$$(gofmt -l .)" || (echo "Files requiring gofmt" && exit 1)
	@test -z "$$(goimports -l .)" || (echo "Files requiring goimports" && exit 1)

test-graph:
	go test ./graphql/... -count=1

test-rest:
	go test ./rest/... -count=1

test-shared:
	go test ./shared/... -count=1

test: test-shared test-rest test-graph

docs-rest:
	swag fmt -d rest
	swag init -d rest -g server.go -o rest/docs --outputTypes yaml,go

docs: docs-rest

build-rest: docs-rest
	go build -o $(REST_EXEC_NAME) ./rest

build-graph:
	go build -o $(GRAPH_EXEC_NAME) ./graphql

build: build-rest build-graph

clean-rest:
	rm -f $(REST_EXEC_NAME) rest/$(REST_EXEC_NAME) $(REST_EXEC_NAME).exe rest/$(REST_EXEC_NAME).exe

clean-graph:
	rm -f $(GRAPH_EXEC_NAME) graphql/$(GRAPH_EXEC_NAME) $(GRAPH_EXEC_NAME).exe graphql/$(GRAPH_EXEC_NAME).exe

clean: clean-rest clean-graph

docker-rest:
	$(DOCKER_RUNNER) build -f rest/Dockerfile -t $(REST_BASE_TAG):$(RELEASE_TAG) .
	$(DOCKER_RUNNER) tag $(REST_BASE_TAG):$(RELEASE_TAG) $(REST_BASE_TAG):latest

docker-graph:
	$(DOCKER_RUNNER) build -f graphql/Dockerfile -t $(GRAPH_BASE_TAG):$(RELEASE_TAG) .
	$(DOCKER_RUNNER) tag $(GRAPH_BASE_TAG):$(RELEASE_TAG) $(GRAPH_BASE_TAG):latest

docker: docker-rest docker-graph

