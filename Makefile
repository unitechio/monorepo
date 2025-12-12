# Makefile (oss-monorepo)
.PHONY: help dev goreman build backend frontend docker-build up down logs clean
.PHONY: build-dev build-prod docker-build-dev docker-build-prod

# environment
ENV ?= development

help:
	@echo "Targets:"
	@echo "  make dev               # run dev via goreman"
	@echo "  make build-dev         # build for development"
	@echo "  make build-prod        # build for production"
	@echo "  make docker-build-dev  # build docker image (dev)"
	@echo "  make docker-build-prod # build docker image (prod)"
	@echo "  make clean             # remove artifacts"

dev: prepare
	@goreman -f Procfile start

goreman: dev

prepare:
	@mkdir -p api/cmd/server/out
	@if [ ! -f api/cmd/server/out/index.html ]; then echo "<html><body><h1>Backend Running (Dev Mode)</h1></body></html>" > api/cmd/server/out/index.html; fi

backend: prepare
	@cd api && go run ./cmd/server

frontend:
	@cd web && npm install && npm run dev

build-prod:
	@echo "Building (Production)..."
	@cd web && npm ci && npm run build
	@echo "Copying assets..."
	@rm -rf api/cmd/server/out
	@cp -r web/out api/cmd/server/out
	@echo "Building backend..."
	@cd api && CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o bin/backend ./cmd/server

build-dev:
	@echo "Building (Development)..."
	@cd web && npm install && npm run build
	@echo "Copying assets..."
	@rm -rf api/cmd/server/out
	@cp -r web/out api/cmd/server/out
	@echo "Building backend..."
	@cd api && go build -o bin/backend ./cmd/server

build: build-prod

docker-build-prod:
	@docker build --build-arg APP_ENV=production -t oss-monorepo:prod .

docker-build-dev:
	@docker build --build-arg APP_ENV=development -t oss-monorepo:dev .

docker-build: docker-build-prod

clean:
	@rm -rf api/bin
	@rm -rf api/cmd/server/out
	@rm -rf web/out
	@rm -rf web/.next

up:
	@docker-compose up -d

down:
	@docker-compose down

logs:
	@docker-compose logs -f
