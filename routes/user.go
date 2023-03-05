package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sing3demons/gin-backend-api/controllers/user"
	"gorm.io/gorm"
)

func InitUserRoutes(router *gin.RouterGroup, db *gorm.DB) {
	userRouter := router.Group("/users")
	userController := user.New(db)

	userRouter.GET("/", userController.GetAll)
	userRouter.GET("/:id", userController.GetById)
	userRouter.POST("/register", userController.Register)
	userRouter.POST("/login", userController.Login)
}
