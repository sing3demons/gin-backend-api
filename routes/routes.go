package routes

import (
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	logMiddleware "github.com/sing3demons/gin-backend-api/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func NewRouter(logger *zap.Logger, db *gorm.DB) *gin.Engine {
	router := gin.Default()
	router.Use(logMiddleware.ZapLogger(logger))
	router.Use(ginzap.RecoveryWithZap(logger, true))
	apiV1 := router.Group("/api/v1")
	InitUserRoutes(apiV1, db)

	return router
}
