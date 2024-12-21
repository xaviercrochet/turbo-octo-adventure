build: build-api build-web

build-api:
	go build -o bin/api ./cmd/api/main.go

build-web:
	go build -o bin/web ./cmd/web/main.go
