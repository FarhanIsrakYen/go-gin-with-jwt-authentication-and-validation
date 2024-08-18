package routes

import (
	"github.com/gin-gonic/gin"
	controller "go-gin-with-jwt-authentication-and-validation/controllers"
	"go-gin-with-jwt-authentication-and-validation/middleware"
)

func ApiRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/users", controller.GetUsers())
	incomingRoutes.GET("/users/:id", controller.GetUser())
}
