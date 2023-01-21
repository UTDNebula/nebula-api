

ifndef DOCKER_RUNNER
	DOCKER_RUNNER = docker
endif

ifndef CONTAINER_HOST
	REGISTRY = localhost:5000
endif

build: server.go
	go build -o nebula-api server.go

clean: nebula-api
	rm nebula-api

docker: Dockerfile server.go
	$(DOCKER_RUNNER) build -t $(REGISTRY)/nebula-api/go-api .
