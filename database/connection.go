package database

import (
	"go-gin-with-jwt-authentication-and-validation/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Connect establishes a connection to the database using Gorm and returns the *gorm.DB instance.
func Connect() (*gorm.DB, error) {
	//	dbHost := os.Getenv("DB_HOST")
	//	dbPort := os.Getenv("DB_PORT")
	//	dbUser := os.Getenv("DB_USER")
	//	dbPassword := os.Getenv("DB_PASSWORD")
	//	dbName := os.Getenv("DB_NAME")

	dsn := "youruser:yourpassword@tcp(127.0.0.1:3306)/yourdbname?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.User{})

	return db, nil
}
