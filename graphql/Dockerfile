# syntax=docker/dockerfile:1

FROM docker.io/golang:1.26 AS builder
WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY ./graphql ./graphql
COPY ./shared ./shared

RUN CGO_ENABLED=0 go build -v -o /build/graph-api ./graphql

FROM gcr.io/distroless/static-debian11
WORKDIR /app
COPY --from=builder /build/graph-api /app/graph-api
ENTRYPOINT [ "/app/graph-api" ]
