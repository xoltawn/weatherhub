DOCKER_COMPOSE=docker compose

run:
	swag init -g cmd/server/main.go
	go run cmd/server/main.go

build:
	go build -o bin/server cmd/server/main.go

test:
	go test ./...

tidy:
	go mod tidy

swag:
	swag init -g cmd/server/main.go


up:
	@echo "Starting containers..."
	$(DOCKER_COMPOSE) up -d
	docker logs weather_hub_api  -f
	
down:
	@echo "Stopping containers..."
	$(DOCKER_COMPOSE) down

config:
	cp .env.example .env

generate:
	go generate ./...

install-tools:
	@echo "Installing dev tools..."
	go install github.com/vektra/mockery/v2@latest
	go install github.com/swaggo/swag/cmd/swag@latest

.PHONY: run build test tidy swag up down
