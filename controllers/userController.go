package controllers

import (
	"github.com/gin-gonic/gin"
	helper "go-gin-with-jwt-authentication-and-validation/helpers"
	"go-gin-with-jwt-authentication-and-validation/models"
	"net/http"
	"time"
)

func HashPassword() {

}

func VerifyPassword() {

}

func GetUsers() gin.HandlerFunc {
	return func(context *gin.Context) {
		userId := context.Param("id")
		if err := helper.MatchUserTypeToUid(context, userId); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var context, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User

	}
}

func GetUser() gin.HandlerFunc {
	return func(context *gin.Context) {

	}
}

func Signup() gin.HandlerFunc {
	return func(context *gin.Context) {

	}
}

func Login() gin.HandlerFunc {
	return func(context *gin.Context) {

	}
}
