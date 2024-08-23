package controllers

import (
	"go-gin-with-jwt-authentication-and-validation/database"
	"go-gin-with-jwt-authentication-and-validation/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct{}

func (ac UserController) GetUser(c *gin.Context) {
    userId, exists := c.Get("userId")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    var user models.User
    if err := database.DB.First(&user, userId).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    c.JSON(http.StatusOK, user)
}

func (ac UserController) UpdateUser(c *gin.Context) {
	var updateUser struct {
        Username        string `json:"username"`
        CurrentPassword string `json:"current_password"`
        NewPassword     string `json:"new_password"`
    }
    
    if err := c.ShouldBindJSON(&updateUser); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }
    
    userId, exists := c.Get("userId")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    
    var user models.User
    if err := database.DB.First(&user, userId).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    
    if updateUser.NewPassword != "" {
        if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(updateUser.CurrentPassword)); err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Current password is incorrect"})
            return
        }
    
        // Hash the new password
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateUser.NewPassword), bcrypt.DefaultCost)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new password"})
            return
        }
        user.Password = string(hashedPassword)
    }
    
    if updateUser.Username != "" {
        user.Username = updateUser.Username
    }
    
    if err := database.DB.Save(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (ac UserController) DeactivateUser(c *gin.Context) {
    username, exists := c.Get("username")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    var user models.User
    database.DB.Where("username = ?", username).First(&user)

    if user.ID == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    database.DB.Model(&models.User{}).Where("username = ?", username).Update("is_active", false)
    c.JSON(http.StatusOK, gin.H{"message": "User deactivated successfully"})
}
