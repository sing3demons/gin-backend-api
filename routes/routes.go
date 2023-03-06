package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	logMiddleware "github.com/sing3demons/gin-backend-api/logger"
	prometheus "github.com/zsais/go-gin-prometheus"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func NewRouter(logger *zap.Logger, db *gorm.DB) *gin.Engine {
	router := gin.Default()
	p := prometheus.NewPrometheus("gin")
	p.Use(router)
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	router.Use(logMiddleware.ZapLogger(logger))
	router.Use(logMiddleware.RecoveryWithZap(logger, true))
	apiV1 := router.Group("/api/v1")
	InitHomeRoutes(apiV1)
	InitUserRoutes(apiV1, db)
	InitBlogRoutes(apiV1, db)

	return router
}
