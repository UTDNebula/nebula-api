

ifndef DOCKER_RUNNER
	DOCKER_RUNNER = docker
endif

ifndef REGISTRY 
	REGISTRY = localhost:5000
endif

build: server.go
	go build -o go-api server.go

clean: nebula-api
	rm nebula-api

RELEASE_TAG=$(shell git rev-parse --short HEAD)
BASE_TAG=$(REGISTRY)/utdnebula/api/go-api
docker: Dockerfile configs/ controllers/ models/ responses/ routes/ server.go
	$(DOCKER_RUNNER) build -t $(BASE_TAG):$(RELEASE_TAG) .
	$(DOCKER_RUNNER) tag $(BASE_TAG):$(RELEASE_TAG) $(BASE_TAG):latest
