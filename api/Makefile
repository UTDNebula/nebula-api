
EXEC_NAME?=go-api
DOCKER_RUNNER?=docker
REGISTRY?=localhost:5000
RELEASE_TAG=$(shell git rev-parse --short HEAD)
BASE_TAG=$(REGISTRY)/utdnebula/api/go-api

build: server.go
	go build -o $(EXEC_NAME) server.go

clean: $(EXEC_NAME)
	rm $(EXEC_NAME)

docker: Dockerfile configs/ controllers/ models/ responses/ routes/ server.go
	$(DOCKER_RUNNER) build -t $(BASE_TAG):$(RELEASE_TAG) .
	$(DOCKER_RUNNER) tag $(BASE_TAG):$(RELEASE_TAG) $(BASE_TAG):latest
