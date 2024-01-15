// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// Swagger docs.
	_ "jastip-app/docs"
	"jastip-app/internal/usecase"
)

// NewRouter -.
// Swagger spec:
// @title       Go Clean Template API
// @version     1.0
// @host        localhost:8080
// @BasePath    /api/v1
func NewRouter(handler *gin.Engine, uc *usecase.UseCase) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "Ok"}) })

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routers
	handler.Group("/api/v1.0")

}
