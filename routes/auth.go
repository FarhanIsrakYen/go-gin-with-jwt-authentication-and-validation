package routes

import (
	"github.com/gin-gonic/gin"
	controller "go-gin-with-jwt-authentication-and-validation/controllers"
)

func Authentication(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("users/signup", controller.Signup())
	incomingRoutes.POST("users/login", controller.Login())
}
