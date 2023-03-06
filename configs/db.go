package configs

import (
	"fmt"
	"os"

	"github.com/sing3demons/gin-backend-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connection() *gorm.DB {
	dsn := os.Getenv("DATABASE_DSN")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Database connected")
	db.AutoMigrate(&models.User{}, &models.Blog{})
	return db
}
