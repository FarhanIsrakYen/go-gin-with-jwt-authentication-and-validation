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
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	router := gin.Default()

    err := routes.SetupRoutes(router)
    if err != nil {
        fmt.Printf("Error setting up routes: %v", err)
    }

	router.Run(":" + port)
}
