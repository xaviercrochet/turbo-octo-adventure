FROM golang:1.23-alpine
RUN apk add git


WORKDIR /app
COPY ../../go.mod ../../go.sum ./
COPY ../../web ./web
COPY ../../pkg ./pkg
COPY ../../cmd/web ./cmd/web

RUN go mod download 
RUN go build -o ./bin/web ./cmd/web/main.go

# stage 2
FROM alpine:3.19

WORKDIR /app

COPY --from=0 /app/bin/web .
