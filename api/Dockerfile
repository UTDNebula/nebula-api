# syntax=docker/dockerfile:1

FROM golang:1.18-alpine

WORKDIR /app
    
COPY go.mod ./
COPY go.sum ./
RUN go mod download

ADD . /app

RUN go build -o /api

EXPOSE 8080

CMD [ "/api" ]