package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sing3demons/gin-backend-api/controllers/blog"
	"github.com/sing3demons/gin-backend-api/middleware"
	"gorm.io/gorm"
)

func InitBlogRoutes(router *gin.RouterGroup, db *gorm.DB) {
	blogRouter := router.Group("/blogs")
	blogController := blog.New(db)
	protect := middleware.AuthorizeJWT()

	blogRouter.GET("/", blogController.GetAll)
	blogRouter.GET("/:id", blogController.GetById)
	blogRouter.POST("/", protect, blogController.Create)
	blogRouter.PATCH("/:id", blogController.Update)
	blogRouter.DELETE("/:id", blogController.Delete)
}
