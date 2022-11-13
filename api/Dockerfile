# syntax=docker/dockerfile:1

FROM golang:1.18-alpine

WORKDIR /app
    
COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -v -o go-api

CMD [ "/app/go-api" ]