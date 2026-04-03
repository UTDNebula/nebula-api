EXEC_NAME?=go-graph
DOCKER_RUNNER?=docker
REGISTRY?=localhost:5001
RELEASE_TAG=$(shell git rev-parse --short HEAD)
BASE_TAG=$(REGISTRY)/utdnebula/api/go-graph

setup:
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install golang.org/x/tools/cmd/goimports@latest

check:
	go mod tidy
	go vet ./...
	staticcheck ./...
	gofmt -w ./..
	goimports -w ./..

generate:
	go get github.com/99designs/gqlgen@latest
	go run github.com/99designs/gqlgen generate
	go mod tidy

test:
	go test ./... -count=1

build: server.go
	go build -o $(EXEC_NAME) server.go

clean: $(EXEC_NAME)
	rm $(EXEC_NAME)

docker: Dockerfile configs/ graph/ gqlgen.yml server.go
	$(DOCKER_RUNNER) build -t $(BASE_TAG):$(RELEASE_TAG) .
	$(DOCKER_RUNNER) tag $(BASE_TAG):$(RELEASE_TAG) $(BASE_TAG):latest
