package route

import (
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unitechio/oss-monorepo/server/internal/config"
	"github.com/unitechio/oss-monorepo/server/pkg/apperr"
	"github.com/unitechio/oss-monorepo/server/pkg/response"
	"github.com/unitechio/oss-monorepo/server/pkg/utils/constants"
	"github.com/unitechio/oss-monorepo/server/web"
)

type Handlers struct {
}

func SetupRoutes(r *gin.Engine, cfg *config.Config, h Handlers) {
	authService := newMockAuthService(cfg)

	r.GET("/health", func(c *gin.Context) {
		response.OK(c, "healthy", gin.H{"service": "api"})
	})

	v1 := r.Group(constants.RequestMappingV1)
	{
		authService.registerRoutes(v1)
		v1.GET("/ping", func(c *gin.Context) {
			response.OK(c, "pong", nil)
		})
	}

	// Serve Frontend Embedded Files
	frontendFS := web.FS()
	staticServer := gin.WrapH(http.FileServer(http.FS(frontendFS)))

	// Catch-all route to serve static files or fallback to index.html for SPA routing
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if path == "" {
			path = "/"
		}

		if path == constants.RequestMappingV1 || len(path) > len(constants.RequestMappingV1) && path[:len(constants.RequestMappingV1)] == constants.RequestMappingV1 {
			response.FailError(c, apperr.NotFound("api route", path))
			return
		}

		// Check if the file exists in the embedded filesystem
		assetPath := "index.html"
		if path != "/" {
			assetPath = path[1:] // remove leading slash
		}
		_, err := fs.Stat(frontendFS, assetPath)
		if err == nil {
			// File exists, serve it
			staticServer(c)
			return
		}

		// File doesn't exist, fallback to index.html for React SPA
		c.Request.URL.Path = "/"
		staticServer(c)
	})
}
