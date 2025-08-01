// main.go
package main

import (
	"sortlynk/database"
	"sortlynk/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database connections
	database.Connect()

	// Initialize Gin router
	r := gin.Default()

	// Apply global middleware
	r.Use(handlers.AuthMiddleware())

	// Public routes (no rate limiting for redirects)
	r.GET("/:code", handlers.RedirectURL)

	// API routes with rate limiting
	api := r.Group("/api/v1")
	api.Use(handlers.RateLimitMiddleware())
	{
		// Authentication routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
		}

		// URL routes
		urls := api.Group("/urls")
		{
			urls.POST("/shorten", handlers.ShortenURL)
			urls.GET("/my", handlers.GetUserURLs)
			urls.GET("/:code/stats", handlers.GetURLStats)
		}
	}

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.Run(":8080")
}
