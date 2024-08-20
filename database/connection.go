package database

import (
	"fmt"
	"go-gin-with-jwt-authentication-and-validation/models"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Connect establishes a connection to the database using Gorm and returns the *gorm.DB instance.
func Connect() () {
	var err error
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dbUser := os.Getenv("DB_USER")
		dbPassword := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",dbUser, dbPassword, dbHost, dbPort, dbName)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to the database:", err)
    }

    err = DB.AutoMigrate(&models.User{})
    if err != nil {
        log.Fatal("Failed to migrate database:", err)
    }
	fmt.Println("Database connection established")
}
