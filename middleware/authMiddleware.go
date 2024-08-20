package middleware

import (
	"go-gin-with-jwt-authentication-and-validation/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
            c.Abort()
            return
        }

        claims := &controllers.Claims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return controllers.jwtKey, nil // Correct reference to jwtKey
		})

        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        c.Set("username", claims.Username)
        c.Set("role", controllers.ROLE_USER)
        c.Next()
    }
}

func AdminOnly(c *gin.Context) {
    role, exists := c.Get("role")
    if !exists || role != "ROLE_ADMIN" {
        c.JSON(http.StatusForbidden, gin.H{"error": "Access Denied"})
        c.Abort()
        return
    }
    c.Next()
}

func UserOnly(c *gin.Context) {
    role, exists := c.Get("role")
    if !exists || role != controllers.ROLE_USER || role != controllers.ROLE_ADMIN {
        c.JSON(http.StatusForbidden, gin.H{"error": "Access Denied"})
        c.Abort()
        return
    }
    c.Next()
}
