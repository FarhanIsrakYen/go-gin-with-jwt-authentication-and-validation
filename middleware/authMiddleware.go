package middleware

import (
	"fmt"
	"go-gin-with-jwt-authentication-and-validation/controllers"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

func JWTAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }

        // Extract token from Authorization header
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenString == authHeader { // No "Bearer " prefix found
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
            c.Abort()
            return
        }

        // Parse the token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            // Validate the signing method
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return jwtKey, nil
        })

        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "details": err.Error()})
            c.Abort()
            return
        }

        // Validate token claims
        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            fmt.Println(claims)
            userId, ok := claims["user_id"].(string)
            if !ok {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
                c.Abort()
                return
            }
            c.Set("userId", userId)
        } else {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

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
