# Makefile (oss-monorepo)
.PHONY: help setup dev goreman build backend frontend docker-build up down logs clean
.PHONY: build-dev build-prod docker-build-dev docker-build-prod

# environment
ENV ?= development

setup:
	@echo "Setting up project..."
	@echo "Installing frontend dependencies..."
	@npm install
	@echo "Installing backend dependencies..."
	@go mod download
	@go mod tidy
	@echo "Creating necessary directories..."
	@mkdir -p api/cmd/server/out
	@mkdir -p tmp
	@mkdir -p out
	@if [ ! -f api/cmd/server/out/index.html ]; then echo "<html><body><h1>Backend Running (Dev Mode)</h1></body></html>" > api/cmd/server/out/index.html; fi
	@echo "Setup completed!"
	@echo ""
	@echo "Note: Make sure you have:"
	@echo "   - PostgreSQL database running (default: localhost:5432)"
	@echo "   - Database created (default: oss_monorepo)"
	@echo "   - Environment variables configured (if needed)"
	@echo ""
	@echo "Run 'make dev' to start development servers"

help:
	@echo "Targets:"
	@echo "  make setup             # 📦 Setup project (install dependencies)"
	@echo "  make dev               # 🚀 Run both backend and frontend (development mode)"
	@echo "  make backend           # Run backend only"
	@echo "  make frontend          # Run frontend only"
	@echo "  make build-dev         # Build for development"
	@echo "  make build-prod        # Build for production"
	@echo "  make docker-build-dev  # Build docker image (dev)"
	@echo "  make docker-build-prod # Build docker image (prod)"
	@echo "  make clean             # Remove artifacts"

dev: prepare
	@echo "🚀 Starting backend and frontend..."
	@goreman -f Procfile start

goreman: dev

prepare:
	@mkdir -p api/cmd/server/out
	@if [ ! -f api/cmd/server/out/index.html ]; then echo "<html><body><h1>Backend Running (Dev Mode)</h1></body></html>" > api/cmd/server/out/index.html; fi

backend: prepare
	@go run ./api/cmd/server

frontend:
	@npm install && npm run dev

build-prod:
	@echo "Building (Production)..."
	@npm ci && npm run build
	@echo "Copying assets..."
	@rm -rf api/cmd/server/out
	@cp -r out api/cmd/server/out
	@echo "Building backend..."
	@cd api && CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o bin/backend ./cmd/server

build-dev:
	@echo "Building (Development)..."
	@npm install && npm run build
	@echo "Copying assets..."
	@rm -rf api/cmd/server/out
	@cp -r out api/cmd/server/out
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
	@rm -rf out
	@rm -rf .next

up:
	@docker-compose up -d

down:
	@docker-compose down

logs:
	@docker-compose logs -f
