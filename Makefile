.PHONY: run generate up down

BACKEND_DIR := backend
INFRA_DIR   := infra

run:
	cd $(BACKEND_DIR) && go run ./cmd/server/*.go

generate:
	cd $(BACKEND_DIR) && go generate ./...

up:
	cd $(INFRA_DIR) && docker compose up -d

down:
	cd $(INFRA_DIR) && docker compose down
