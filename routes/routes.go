package routes

import (
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func NewRouter(logger *zap.Logger, db *gorm.DB) *gin.Engine {
	router := gin.Default()

	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(logger, true))

	apiV1 := router.Group("/api/v1")
	InitUserRoutes(apiV1, db)

	return router
}
