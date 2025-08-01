package handlers

import (
	"context"
	"fmt"
	"net/http"
	"sortlynk/database"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("authenticated", true)
		c.Next()
	}
}

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAuthenticated, exists := c.Get("authenticated")
		if exists && isAuthenticated.(bool) {
			// Authenticated users have higher limits
			c.Next()
			return
		}

		clientIP := c.ClientIP()
		key := fmt.Sprintf("rate_limit:%s", clientIP)

		ctx := context.Background()
		count, err := database.Redis.Get(ctx, key).Int()
		if err != nil && err.Error() != "redis: nil" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Rate limiting error"})
			c.Abort()
			return
		}

		if count >= 200 {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Unauthenticated users can make 2 requests per day.",
			})
			c.Abort()
			return
		}

		// Increment counter
		database.Redis.Incr(ctx, key)
		database.Redis.Expire(ctx, key, 24*time.Hour)

		c.Next()
	}
}
