# Project Structure

```
NEBULA-API
├───api                 - API
│   ├───.env            - API developer environment variables (MONGODB_URI)
│   ├───server.go       - Server (go run server.go)
│   ├───configs         - MongoDB configuration
│   ├───controllers     - Route controllers
│   ├───models          - Schema models
│   ├───responses       - Route responses
│   ├───routes          - Define routes
│   └schema
│      ├───objects.go      - Object schemas
│      └───requirements.go - Requirement schemas
│───docs
│   ├───schemas         - Database schemas documentation
│   └───index.md        - Project structure and usage
└───README.md           - Project overview
```

# Building

## Standalone Executable

The easiest way to build the project is to use the Makefile. To build this, simply run `make build`. To clean the executable, run `make clean`. Note that this will fail if the build executable does not exist. If you do not have make, you can look at the build and clean targets and run them accordingly in the Makefile.

For ease of building in Windows environments, a `build.bat` alternative is provided. Note that unlike the makefile, this does NOT build a docker image.

## Docker

To build the docker image for the API, run `make docker`. This will run the build command on any docker runner (default is docker) and tag it accordingly:

```
(REGISTRY)/utdnebula/api/go-api:(first six of git hash)
```

and

```
(REGISTRY)/utdnebula/api/go-api:latest
```

By default, REGISTRY is set to `localhost:5000`. You can set this as a different environment variable if you want to upload the docker image to a different registry.

# Testing

## Linting

To run `staticcheck` refer to [this URL](https://staticcheck.io/docs/getting-started/). Run `make check` to run `go vet` and `staticcheck` on the api project. In order to a PR to be approved, this make target (or its relevant commands) should run successfully to ensure use of good practices.

## Unit/Integration Testing (TODO)

# Running

## Standalone

To connect to your Mongo database, in /api/.env: set MONGODB_URI accordingly:

```
MONGODB_URI=<insert_connection_string_here>
```

You may also specify a different port (if 8080 is not desired) in the same .env file by setting PORT as well.

Then, after running `make build`, run the `go-api` executable. Alternatively, you can run `server.go` directly with `go run server.go`. The server will begin serving the API on "/".

View the swagger-ui at [http://localhost:8080/swagger/index.html] port may be different if alternatively set

## Docker

After building the image, create the .env file just as described earlier. Next, run the following docker command:

`docker run -d -p "8080:8080" -v "./env:/app/.env" localhost:5000/utdnebula/api/go-api:latest`

Change `localhost:5000` to a different registry if you customized this during the docker build.
