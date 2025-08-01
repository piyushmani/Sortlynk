package models

import (
	"time"

	"gorm.io/gorm"
)

type URL struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	ShortCode   string         `json:"short_code" gorm:"unique;not null;index"`
	OriginalURL string         `json:"original_url" gorm:"not null"`
	UserID      *uint          `json:"user_id,omitempty"`
	User        *User          `json:"user,omitempty"`
	ClickCount  int64          `json:"click_count" gorm:"default:0"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}