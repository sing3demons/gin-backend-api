package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitHomeRoutes(r *gin.RouterGroup) {
	router := r.Group("/")
	router.GET("", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"API VERSION": "1.0.0",
		})
	})
}
