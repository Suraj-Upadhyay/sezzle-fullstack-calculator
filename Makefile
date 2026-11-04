.PHONY: install dev-backend dev-frontend test test-backend test-frontend \
	coverage lint build docker-build docker-run

IMAGE ?= calculator

install:
	cd frontend && npm ci
	cd backend && go mod download

dev-backend:
	cd backend && go run ./cmd/server

dev-frontend:
	cd frontend && npm run dev

test: test-backend test-frontend

test-backend:
	cd backend && go test -race ./...

test-frontend:
	cd frontend && npm test

coverage:
	mkdir -p coverage
	cd backend && go test -coverpkg=./... -coverprofile=coverage.out ./... \
		&& go tool cover -func=coverage.out | tail -1 \
		&& go tool cover -html=coverage.out -o ../coverage/backend.html
	cd frontend && npm run coverage
	@echo "Reports: coverage/backend.html, frontend/coverage/index.html"

lint:
	cd backend && go vet ./... && golangci-lint run ./...
	cd frontend && npm run lint && npm run format:check && npm run typecheck

build:
	cd frontend && npm run build
	cd backend && go build -o ../dist/server ./cmd/server

docker-build:
	docker build -t $(IMAGE) .

docker-run: docker-build
	docker run --rm -p 8080:8080 $(IMAGE)
