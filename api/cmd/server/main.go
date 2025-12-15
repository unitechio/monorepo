package main

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/unitechio/oss-backend/internal/infrastructure/database"
	"github.com/unitechio/oss-backend/internal/repository"
	"github.com/unitechio/oss-backend/internal/service"
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
	// Initialize Database
	database.InitDB()

	// Initialize Layers
	pageRepo := repository.NewPageRepository()
	service.InitPageService(pageRepo)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
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

	r.Run(":8080")
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
