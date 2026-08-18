# syntax=docker/dockerfile:1

FROM docker.io/golang:1.25 AS builder
WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY ./rest ./rest

RUN CGO_ENABLED=0 go build -v -o /build/go-api ./rest

FROM gcr.io/distroless/static-debian11
WORKDIR /app
COPY --from=builder /build/go-api /app/go-api
ENTRYPOINT [ "/app/go-api" ]

