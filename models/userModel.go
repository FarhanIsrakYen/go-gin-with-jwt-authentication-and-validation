package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
    gorm.Model
    Username string `gorm:"unique;not null"`
    Password string `gorm:"not null"`
    Email    string `gorm:"unique;not null"`
    UserRole string `gorm:"not null"` // ROLE_ADMIN or ROLE_USER
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
