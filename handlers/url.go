package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/nelsonandreass/url-shortener/db"
	"github.com/nelsonandreass/url-shortener/models"
	"github.com/nelsonandreass/url-shortener/response"
	"github.com/teris-io/shortid"
)

func ShortenURL(c *gin.Context) {
	var input models.URL

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
	}
	input.ShortURL = shortid.MustGenerate()
	db.DB.Create(&input)
	response := response.UrlResponse{
		OriginalURL: input.OriginalURL,
		ShortURL:    os.Getenv("BASE_URL") + ":" + os.Getenv("BASE_URL_PORT") + "/" + input.ShortURL,
		Hits:        input.Hits,
	}
	c.JSON(http.StatusOK, response)
}

func RedirectURL(c *gin.Context) {
	var url models.URL
	if err := db.DB.Where("short_url = ? ", c.Param("short_url")).First(&url).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
	}
	url.Hits++
	db.DB.Save(&url)
	c.Redirect(http.StatusMovedPermanently, url.OriginalURL)
}

func GetHits(c *gin.Context) {
	var request models.URL

	if err := c.ShouldBindBodyWithJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
	}

	if err := db.DB.Where("short_url = ? ", request.ShortURL).First(&request).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
	}
	response := response.UrlResponse{
		OriginalURL: request.OriginalURL,
		ShortURL:    os.Getenv("BASE_URL") + request.ShortURL,
		Hits:        request.Hits,
	}
	c.JSON(http.StatusOK, response)
}
