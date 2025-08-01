package handlers

import (
	"context"
	"net/http"
	"net/url"
	"sortlynk/database"
	"sortlynk/models"
	"sortlynk/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ShortenRequest struct {
	URL string `json:"url" binding:"required,url"`
}

func ShortenURL(c *gin.Context) {
	var req ShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := url.ParseRequestURI(req.URL); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL format"})
		return
	}

	var urlRecord models.URL
	var err error
	for i := 0; i < 5; i++ {
		shortCode := utils.GenerateShortUrl(req.URL)
		urlRecord = models.URL{
			ShortCode:   shortCode,
			OriginalURL: req.URL,
		}

		if userID, exists := c.Get("user_id"); exists {
			userIDUint := userID.(uint)
			urlRecord.UserID = &userIDUint
		}

		err = database.DB.Create(&urlRecord).Error
		if err == nil {
			break // Success
		}

		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			continue // Collision, try again
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create short URL"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create short URL after multiple attempts"})
		return
	}

	ctx := context.Background()
	database.Redis.Set(ctx, urlRecord.ShortCode, req.URL, 24*time.Hour)

	c.JSON(http.StatusCreated, gin.H{
		"short_code":   urlRecord.ShortCode,
		"short_url":    "http://localhost:8080/" + urlRecord.ShortCode,
		"original_url": req.URL,
	})
}

func RedirectURL(c *gin.Context) {
	shortCode := c.Param("code")

	ctx := context.Background()

	originalURL, err := database.Redis.Get(ctx, shortCode).Result()
	if err == nil {
		go func() {
			database.DB.Model(&models.URL{}).Where("short_code = ?", shortCode).
				UpdateColumn("click_count", gorm.Expr("click_count + ?", 1))
		}()

		c.Redirect(http.StatusMovedPermanently, originalURL)
		return
	}

	var urlRecord models.URL
	if err := database.DB.Where("short_code = ?", shortCode).First(&urlRecord).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
		return
	}

	database.DB.Model(&urlRecord).UpdateColumn("click_count", gorm.Expr("click_count + ?", 1))
	database.Redis.Set(ctx, shortCode, urlRecord.OriginalURL, 24*time.Hour)

	c.Redirect(http.StatusMovedPermanently, urlRecord.OriginalURL)
}

func GetUserURLs(c *gin.Context) {
	userID := c.GetUint("user_id")

	var urls []models.URL
	if err := database.DB.Where("user_id = ?", userID).Find(&urls).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch URLs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"urls": urls})
}

func GetURLStats(c *gin.Context) {
	shortCode := c.Param("code")
	userID := c.GetUint("user_id")

	var urlRecord models.URL
	query := database.DB.Where("short_code = ?", shortCode)

	if userID != 0 {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&urlRecord).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found or access denied"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"short_code":   urlRecord.ShortCode,
		"original_url": urlRecord.OriginalURL,
		"click_count":  urlRecord.ClickCount,
		"created_at":   urlRecord.CreatedAt,
	})
}
