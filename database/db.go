package database

import (
	"fmt"
	"log"
	"sortlynk/config"
	"sortlynk/models"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB    *gorm.DB
	Redis *redis.Client
)

func Connect() {
	cfg := config.Load()

	var err error
	DB, err = gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = DB.AutoMigrate(&models.User{}, &models.URL{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	Redis = redis.NewClient(&redis.Options{
		Addr:     cfg.GetRedisAddr(),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	fmt.Println("Database and Redis connected successfully")
}
