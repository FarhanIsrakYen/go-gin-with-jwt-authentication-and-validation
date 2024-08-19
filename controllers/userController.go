package controllers

import (
	"go-gin-with-jwt-authentication-and-validation/database"
	"go-gin-with-jwt-authentication-and-validation/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword() {

}

func VerifyPassword() {

}

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		//userId := c.Param("id")
		//if err := helper.MatchUserTypeToUid(c, userId); err != nil {
		//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		//	return
		//}
		//
		//var context, cancel = c.WithTimeout(c.Background(), 100*time.Second)
		//
		//var user models.User

	}
}

func GetUser(context *gin.Context) {
	userId := context.Param("id")
	id, err := strconv.ParseUint(userId, 10, 32)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := GetUserById(uint(id))
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	}

	context.JSON(http.StatusOK, user)
}

func GetUserById(id uint) (*models.User, error) {
	var user models.User

	result := database.DB.First(&user, id)

	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func Signup(context *gin.Context) {
	var body struct {
		FirstName    string   
		LastName     string   
		Password     string   
		Email        string   
		Phone        string
	}
	if context.Bind(&body) != nil {
		context.JSON(http.StatusBadRequest, gin.H {
			"error": "Failed to read body",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}
	hashStr := string(hash)

	user := models.User{
		Email: &body.Email,
		Password: &hashStr,
		FirstName: &body.FirstName,
		LastName: &body.LastName,
		Phone: &body.Phone,
	}

	result := database.DB.Create(&user)

	if result.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create account",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Account created successfully",
	})
}

func Login() gin.HandlerFunc {
	return func(context *gin.Context) {

	}
}
