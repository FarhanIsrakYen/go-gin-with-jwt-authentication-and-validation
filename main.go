package main

import (
	"fmt"
	"go-gin-with-jwt-authentication-and-validation/config"
	"go-gin-with-jwt-authentication-and-validation/database"
	"go-gin-with-jwt-authentication-and-validation/routes"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func init()  {
	config.LoadEnvVariables()
	database.Connect()
	database.SyncDatabase()
}

func main() {
	db, err := database.Connect()
	if err != nil {
		fmt.Println("failed to connect to database: %v", err)
	}

	// Test the connection
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println("failed to get database instance: %v", err)
	}
	defer sqlDB.Close()

	err = sqlDB.Ping()
	if err != nil {
		fmt.Println("failed to ping database: %v", err)
	}

	fmt.Println("Database connection established")

	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routes.Authentication(router)
	routes.ApiRoutes(router)

	router.GET("/api-1", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success": "Access granted for api-1",
		})
	})

	router.GET("/api-2", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success": "Access granted for api-2",
		})
	})

	err = router.Run(":" + port)
	if err != nil {
		return
	}
}
