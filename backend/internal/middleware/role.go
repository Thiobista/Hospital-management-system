package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// RoleMiddleware checks if user has required role
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		tokenString := strings.Replace(header, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		userRole, ok := claims["role"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid role in token"})
			c.Abort()
			return
		}

		// Check if user role is in allowed roles
		allowed := false
		for _, role := range allowedRoles {
			if userRole == role {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		// Store user info in context
		userID, _ := claims["id"].(float64)
		c.Set("userID", uint(userID))
		c.Set("userRole", userRole)

		c.Next()
	}
}

// AdminOnly middleware - only admin can access
func AdminOnly() gin.HandlerFunc {
	return RoleMiddleware("admin")
}

// DoctorOnly middleware - only doctor can access
func DoctorOnly() gin.HandlerFunc {
	return RoleMiddleware("doctor")
}

// ReceptionistOnly middleware - only receptionist can access
func ReceptionistOnly() gin.HandlerFunc {
	return RoleMiddleware("receptionist")
}

// AdminOrDoctor middleware - admin or doctor can access
func AdminOrDoctor() gin.HandlerFunc {
	return RoleMiddleware("admin", "doctor")
}

// AdminOrReceptionist middleware - admin or receptionist can access
func AdminOrReceptionist() gin.HandlerFunc {
	return RoleMiddleware("admin", "receptionist")
}
