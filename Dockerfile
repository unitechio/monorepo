# Stage 1: Build Frontend
FROM node:20-alpine AS frontend-builder
ARG APP_ENV=production
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
# Copy the correct env file to .env.production for Next.js build
COPY .env.${APP_ENV} .env.production
RUN npm run build

# Stage 2: Build Backend
FROM golang:1.24-alpine AS backend-builder
WORKDIR /app/api
COPY api/go.mod api/go.sum ./
RUN go mod download
COPY api/ .
# Copy frontend assets from Stage 1
COPY --from=frontend-builder /app/out ./cmd/server/out
RUN CGO_ENABLED=0 GOOS=linux go build -o /backend ./cmd/server

# Stage 3: Final Image
FROM alpine:latest
WORKDIR /root/
COPY --from=backend-builder /backend .
EXPOSE 8080
ENTRYPOINT ["./backend"]
