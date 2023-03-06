package routes

import (
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	logMiddleware "github.com/sing3demons/gin-backend-api/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func NewRouter(logger *zap.Logger, db *gorm.DB) *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	router.Use(logMiddleware.ZapLogger(logger))
	router.Use(ginzap.RecoveryWithZap(logger, true))
	apiV1 := router.Group("/api/v1")
	InitUserRoutes(apiV1, db)
	InitBlogRoutes(apiV1, db)

	return router
}
