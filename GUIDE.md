# OSS Monorepo Guide

This repository contains a full-stack application with a Go backend and a Next.js frontend. The build system has been configured to embed the frontend into the backend binary for easy single-file deployment.

## Prerequisites

- **Go**: 1.21+
- **Node.js**: 20+
- **Docker**: (Optional, for containerized run)
- **Make**: (Optional, for simplified commands)

## Development (Local)

To run the Frontend and Backend separately in development mode (with hot-reload):

```bash
make dev
# OR
make goreman
```
*   **Backend**: http://localhost:8080
*   **Frontend**: http://localhost:3000

## Architecture: Embedded Frontend

For production, the Frontend is built as a static site and **embedded** directly into the Go binary. This means the Backend serves the Frontend assets.

1.  **Next.js** builds with `output: 'export'` -> `web-monorepo/out`.
2.  **Go** uses `//go:embed` to include `web-monorepo/out` into the binary.
3.  **Gin** serves these files and handles SPA routing (redirects 404s to `index.html` for client-side routing).

## Environments

The project supports `development` (dev) and `production` (prod) environments.

- **Development**:
    - `.env.development` is loaded.
    - Build includes sourcemaps and debug info.
    - Go binary is built without stripping.
- **Production**:
    - `.env.production` is loaded.
    - Build is optimized (Next.js production build).
    - Go binary is stripped (`-s -w`).

### Build Commands

| Environment | Build Command | Docker Command |
| :--- | :--- | :--- |
| **Development** | `make build-dev` | `make docker-build-dev` |
| **Production** | `make build-prod` | `make docker-build-prod` |

## Building (Local Production)

To build the combined binary locally for production:

```bash
make build-prod
```
To run the binary:
```bash
./api/bin/backend
```

## Docker (Production Container)

To build the production image:

```bash
make docker-build-prod
```

Running it:

```bash
make up
```

Stop the container:
```bash
make down
```

## Troubleshooting

- **CORS Errors in Dev**: Ensure `api/cmd/server/main.go` has `web-monorepo`'s URL in `AllowOrigins` (default: `http://localhost:3000`).
- **404 on Refresh**: The Go server is configured with a fallback to `index.html`. If this fails, check the `r.NoRoute` handler in `main.go`.
