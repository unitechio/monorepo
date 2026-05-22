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
	@mkdir -p server/web/dist
	@mkdir -p tmp
	@if [ ! -f server/web/dist/index.html ]; then echo "<html><body><h1>Backend Running (Dev Mode)</h1></body></html>" > server/web/dist/index.html; fi
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
	@mkdir -p server/web/dist
	@if [ ! -f server/web/dist/index.html ]; then echo "<html><body><h1>Backend Running (Dev Mode)</h1></body></html>" > server/web/dist/index.html; fi

backend: prepare
	@go run ./server/cmd/server

frontend:
	@npm install && npm run dev

build-prod:
	@echo "Building (Production)..."
	@npm ci && npm run embed:web:production
	@echo "Building backend..."
	@cd server && CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o bin/backend ./cmd/server

build-dev:
	@echo "Building (Development)..."
	@npm install && npm run embed:web:development
	@echo "Building backend..."
	@cd server && go build -o bin/backend ./cmd/server

build: build-prod

docker-build-prod:
	@docker build --build-arg APP_ENV=production -t oss-monorepo:prod .

docker-build-dev:
	@docker build --build-arg APP_ENV=development -t oss-monorepo:dev .

docker-build: docker-build-prod

clean:
	@rm -rf server/bin
	@rm -rf server/web/dist
	@rm -rf dist

up:
	@docker-compose up -d

down:
	@docker-compose down

logs:
	@docker-compose logs -f
