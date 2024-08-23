package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type Claims struct {
    Username string `json:"username"`
    UserID string `json:"user_id"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}

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
        if tokenString == authHeader { 
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
            c.Abort()
            return
        }

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
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

        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            userId, ok := claims["user_id"].(string)
            if !ok {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
                c.Abort()
                return
            }
            c.Set("userId", userId)
            c.Set("username", claims["username"].(string))
        } else {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }
        
        c.Next()
    }
}

func AuthorizeRole(allowedRoles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
            c.Abort()
            return
        }

        tokenString = strings.TrimPrefix(tokenString, "Bearer ")

        token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })
        claims, ok := token.Claims.(*Claims);

        if ok && token.Valid {
            userRole := claims.Role
            for _, role := range allowedRoles {
                if userRole == role {
                    c.Next()
                    return
                }
            }
            c.JSON(http.StatusForbidden, gin.H{"error": "You don't have access to this resource"})
            c.Abort()
            return
        } else {
            c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
            c.Abort()
            return
        }
    }
}
