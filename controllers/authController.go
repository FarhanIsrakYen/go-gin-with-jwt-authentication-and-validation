package controllers

import (
	"go-gin-with-jwt-authentication-and-validation/config"
	"go-gin-with-jwt-authentication-and-validation/database"
	"go-gin-with-jwt-authentication-and-validation/models"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))
var validate = validator.New()

type Credentials struct {
    Username string `json:"username" validate:"required"`
    Password string `json:"password" validate:"required,min=8"`
}

type Claims struct {
    Username string `json:"username"`
    UserID string `json:"user_id"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}

type AuthController struct{}

func (ac AuthController) SignUp(c *gin.Context) {
    var credentials Credentials
    if err := c.BindJSON(&credentials); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    if err := validate.Struct(&credentials); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }

    user := models.User{
        Username: credentials.Username,
        Password: string(hashedPassword),
        Email:    credentials.Username,
        Role: config.ROLE_USER,
    }

    result := database.DB.Create(&user)
    if result.Error != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Signed up successfully"})
}

func (ac AuthController) Login(c *gin.Context) {
    var credentials Credentials
    if err := c.BindJSON(&credentials); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    if err := validate.Struct(&credentials); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var user models.User
    database.DB.Where("username = ? AND is_active = ?", credentials.Username, true)

    if user.ID == 0 || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)) != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &Claims{
        Username: user.Username,
        UserID: strconv.FormatUint(uint64(user.ID), 10),
        Role:     user.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": tokenString})
}