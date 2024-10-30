package middleware

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	username = "admin"
	password = "password"
)

// BasicAuth is a middleware for Basic Authentication
func BasicAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Check if the header starts with "Basic "
		if !strings.HasPrefix(authHeader, "Basic ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization type must be Basic"})
			c.Abort()
			return
		}

		// Decode the base64 part of the Authorization header
		payload := strings.TrimPrefix(authHeader, "Basic ")

		decoded, err := base64.StdEncoding.DecodeString(payload)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token"})
			c.Abort()
			return
		}

		// Split the decoded string into username and password
		parts := strings.SplitN(string(decoded), ":", 2)
		if len(parts) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		// Validate the username and password
		if parts[0] != username || parts[1] != password {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			c.Abort()
			return
		}

		// If everything is fine, proceed to the next handler
		c.Next()
	}
}
