package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/unitechio/oss-monorepo/api/internal/config"
	"github.com/unitechio/oss-monorepo/api/internal/infrastructure/database"
	"github.com/unitechio/oss-monorepo/api/internal/repository"
	"github.com/unitechio/oss-monorepo/api/internal/service"
)

//go:embed out/*
var staticFiles embed.FS

type Node struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	Address string `json:"address"`
}

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize Database
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize Layers
	pageRepo := repository.NewPageRepository()
	service.InitPageService(pageRepo)

	// Set Gin mode based on environment
	if cfg.Server.Host == "0.0.0.0" && cfg.Server.Port == 8080 {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Use CORS config from config
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORS.AllowedOrigins,
		AllowMethods:     cfg.CORS.AllowedMethods,
		AllowHeaders:     cfg.CORS.AllowedHeaders,
		AllowCredentials: cfg.CORS.AllowCredentials,
	}))

	// API Routes
	api := r.Group("/api/v1")
	{
		api.GET("/nodes", getNodes)
		api.GET("/pages/:slug", getPage)
	}

	r.GET("/healthz", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	// Static File Serving
	files, err := fs.Sub(staticFiles, "out")
	if err != nil {
		panic(err)
	}
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api") {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}

		// Check if file exists in the embedded FS
		_, err := files.Open(strings.TrimPrefix(path, "/"))
		if err == nil {
			c.FileFromFS(path, http.FS(files))
		} else {
			// Fallback to index.html for SPA
			c.FileFromFS("/", http.FS(files))
		}
	})

	// Start server with config
	serverAddr := cfg.Server.GetServerAddr()
	fmt.Printf("Server starting on %s\n", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getNodes(c *gin.Context) {
	sample := []Node{
		{ID: "n-1", Name: "node-alpha", Status: "ready", Address: "10.0.0.10"},
		{ID: "n-2", Name: "node-beta", Status: "offline", Address: "10.0.0.11"},
		{ID: "n-3", Name: "node-gamma", Status: "ready", Address: "10.0.0.12"},
	}

	c.JSON(http.StatusOK, gin.H{"items": sample})
}

func getPage(c *gin.Context) {
	slug := c.Param("slug")
	lang := c.DefaultQuery("lang", "en") // Default to English

	page, err := service.GetPageBySlug(slug, lang)
	if err != nil {
		// Differentiate between 404 and 500 if possible, for now just 500 or 404
		// If using sql.ErrNoRows, we returns 404.
		// Service currently returns err if GetPageBySlug fails.
		// Repo returns generic error.
		// Ideally we check error type.
		if strings.Contains(err.Error(), "no rows") { // Crude check
			c.JSON(http.StatusNotFound, gin.H{"error": "page not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if page == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "page not found"})
		return
	}

	c.JSON(http.StatusOK, page)
}
