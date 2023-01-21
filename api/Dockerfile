# syntax=docker/dockerfile:1

FROM docker.io/golang:1.18-alpine AS builder
WORKDIR /build

COPY ./configs ./configs
COPY ./controllers ./controllers
COPY ./models ./models
COPY ./responses ./responses
COPY ./routes ./routes
COPY go.mod go.sum server.go .

RUN CGO_ENABLED=0 go build -v -o go-api
RUN touch .env

FROM gcr.io/distroless/static-debian11
WORKDIR /app
COPY --from=builder /build/go-api /app/go-api
COPY --from=builder /build/.env /app/.env
ENTRYPOINT [ "/app/go-api" ]
