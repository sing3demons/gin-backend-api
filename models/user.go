package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey;autoIncrement;not null"`
	Fullname  string         `json:"fullname" gorm:"type:varchar(255);not null"`
	Email     string         `json:"email" gorm:"type:varchar(255);not null;unique"`
	Password  string         `json:"-" gorm:"type:varchar(255);not null"`
	IsAdmin   bool           `json:"is_admin" gorm:"type:bool;default:false"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
