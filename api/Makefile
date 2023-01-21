

ifndef DOCKER_RUNNER
	DOCKER_RUNNER = docker
endif

ifndef CONTAINER_HOST
	REGISTRY = localhost:5000
endif

build: server.go
	go build -o go-api server.go

clean: nebula-api
	rm nebula-api

docker: Dockerfile configs/ controllers/ models/ responses/ routes/ server.go
	$(DOCKER_RUNNER) build -t $(REGISTRY)/nebula-api/go-api .
