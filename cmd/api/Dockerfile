FROM golang:1.23-alpine
RUN apk add git


WORKDIR /app
COPY ../../go.mod ../../go.sum ./
COPY ../../api ./api
COPY ../../pkg ./pkg
COPY ../../cmd/api ./cmd/api

RUN go mod download 
RUN go build -o ./bin/api ./cmd/api/main.go

# stage 2
FROM alpine:3.19

WORKDIR /app

COPY --from=0 /app/bin/api .
