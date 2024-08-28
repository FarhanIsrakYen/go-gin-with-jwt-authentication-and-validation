package controllers

import (
	"go-gin-with-jwt-authentication-and-validation/database"
	"go-gin-with-jwt-authentication-and-validation/models"
	"go-gin-with-jwt-authentication-and-validation/repository"
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

    user, err := repository.NewUserRepository(database.DB).FindByUserID(userId)
	if err != nil || !(user.IsActive) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
        return
	}

    c.JSON(http.StatusOK, user)
}

func (ac UserController) UpdateUser(c *gin.Context) {
	var updateUser models.UpdateUser
    
    if err := c.ShouldBindJSON(&updateUser); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }
    
    userId, exists := c.Get("userId")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    
    user, err := repository.NewUserRepository(database.DB).FindByUserID(userId)
	if err != nil || !(user.IsActive) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
        return
	}
    
    if updateUser.NewPassword != "" {
        if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(updateUser.CurrentPassword)); err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Current password is incorrect"})
            return
        }
    
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
    
	err = repository.NewUserRepository(database.DB).UpdateUser(user)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update user"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (ac UserController) DeactivateUser(c *gin.Context) {
    userId, exists := c.Get("userId")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    user, err := repository.NewUserRepository(database.DB).FindByUserID(userId)
	if err != nil || !(user.IsActive) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
        return
	}

    repository.NewUserRepository(database.DB).DeactivateUser(user.ID)
    c.JSON(http.StatusOK, gin.H{"message": "User deactivated successfully"})
}
