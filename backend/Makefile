.PHONY: dev build test migrate lint tidy

dev:
	go run ./cmd/server

build:
	go build -o bin/server ./cmd/server

test:
	go test ./...

migrate:
	@echo "Run the SQL in internal/db/migrations/001_layer1.sql in Supabase SQL editor"

lint:
	golangci-lint run

tidy:
	go mod tidy
