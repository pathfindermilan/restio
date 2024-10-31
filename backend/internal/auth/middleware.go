package auth

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtService JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Println("Missing Authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Println("Malformed Authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		token, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			log.Println("Token validation error:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*jwtCustomClaim); ok && token.Valid {
			log.Printf("Middleware: User ID extracted: %v\n", claims.UserID)
			c.Set("userID", claims.UserID)
		} else {
			log.Println("Invalid token claims")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		c.Next()
	}
}
