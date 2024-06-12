package router

import (
	"github.com/gin-gonic/gin"
	"github.com/nelsonandreass/url-shortener/handlers"
	"github.com/nelsonandreass/url-shortener/middleware"
	"github.com/nelsonandreass/url-shortener/models"
)

func SetupRoutes(r *gin.Engine, limiter *models.RateLimiter) {

	r.POST("/shorten", middleware.AuthMiddleware(), handlers.ShortenURL)
	r.GET("/:short_url", handlers.RedirectURL)
	r.POST("/get-hit-count", middleware.AuthMiddleware(), middleware.RateLimiterMiddleware(limiter), handlers.GetHits)

	r.POST("/login", handlers.Login)
	r.POST("/register", handlers.Register)
}
