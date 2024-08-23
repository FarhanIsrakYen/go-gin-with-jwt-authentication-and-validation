package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
    gorm.Model
    Username string `json:"username" gorm:"unique;not null"`
    Password string `json:"password" gorm:"not null"`
    Email    string `json:"email" gorm:"unique;not null"`
    Role string `json:"role" gorm:"not null"` // ROLE_ADMIN or ROLE_USER
    IsActive  bool   `json:"is_active" gorm:"default:true"`
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
