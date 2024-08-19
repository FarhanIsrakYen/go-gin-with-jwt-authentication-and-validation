package database

import "go-gin-with-jwt-authentication-and-validation/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}