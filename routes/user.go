package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sing3demons/gin-backend-api/controllers/user"
	"github.com/sing3demons/gin-backend-api/middleware"
	"gorm.io/gorm"
)

func InitUserRoutes(router *gin.RouterGroup, db *gorm.DB) {
	userRouter := router.Group("/users")
	userController := user.New(db)
	protect := middleware.AuthorizeJWT()

	userRouter.GET("/", userController.GetAll)
	userRouter.GET("/:id", userController.GetById)
	userRouter.POST("/register", userController.Register)
	userRouter.POST("/login", userController.Login)
	userRouter.GET("profile",protect, userController.GetProfile)
}
