package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/nelsonandreass/url-shortener/db"
	"github.com/nelsonandreass/url-shortener/models"
	"github.com/nelsonandreass/url-shortener/router"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	redisAddr := os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	ctx := context.Background()
	_, err = client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	limiter := models.NewRateLimiter(client, 10, time.Minute)

	r := gin.Default()
	db.ConnectDatabase()
	router.SetupRoutes(r, limiter)

	r.Run()
}
