package bootstrap

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/unitechio/oss-monorepo/server/internal/config"
	"github.com/unitechio/oss-monorepo/server/internal/http/middleware"
	"github.com/unitechio/oss-monorepo/server/internal/http/route"
	"github.com/unitechio/oss-monorepo/server/internal/infrastructure/database"
)

type Application struct {
	Logger *slog.Logger
	Server *http.Server
}

func BuildApplication(cfg *config.Config) (*Application, error) {
	logger := newLogger(cfg)

	if cfg.Database.Enabled {
		_, err := database.InitDatabases(cfg.Database)
		if err != nil {
			return nil, fmt.Errorf("connect database: %w", err)
		}
	} else {
		logger.Info("database integration disabled", "hint", "set DB_ENABLED=true to require PostgreSQL at startup")
	}

	handlers := route.Handlers{}

	router := newRouter(cfg, logger, handlers)
	server := &http.Server{
		Addr:         resolveAddress(cfg),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	return &Application{Logger: logger, Server: server}, nil
}

func newLogger(cfg *config.Config) *slog.Logger {
	var handler slog.Handler
	if cfg.Logging.Format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	}
	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}

func newRouter(cfg *config.Config, logger *slog.Logger, handlers route.Handlers) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(middleware.RequestID())
	r.Use(middleware.Logger(logger))
	r.Use(middleware.ErrorHandler(logger))
	r.Use(middleware.Recovery(logger))
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.CorsMiddleware(cfg.CORS))
	if cfg.RateLimit.Requests > 0 {
		r.Use(middleware.RateLimit(cfg.RateLimit.Requests))
	}

	route.SetupRoutes(r, cfg, handlers)

	return r
}

func resolveAddress(cfg *config.Config) string {
	port := cfg.Server.Port
	if port == 0 {
		port = 8080
	}
	if cfg.Server.Host != "" && cfg.Server.Host != "0.0.0.0" {
		return fmt.Sprintf("%s:%d", cfg.Server.Host, port)
	}
	return fmt.Sprintf(":%d", port)
}
