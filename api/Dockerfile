# syntax=docker/dockerfile:1

FROM docker.io/golang:1.23 AS builder
WORKDIR /build

COPY ./common ./common
COPY ./configs ./configs
COPY ./controllers ./controllers
COPY ./docs ./docs
COPY ./responses ./responses
COPY ./routes ./routes
COPY ./schema ./schema
COPY go.mod go.sum server.go ./

RUN CGO_ENABLED=0 go build -v -o go-api

FROM gcr.io/distroless/static-debian11
WORKDIR /app
COPY --from=builder /build/go-api /app/go-api
ENTRYPOINT [ "/app/go-api" ]
