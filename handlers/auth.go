package handlers

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/nelsonandreass/url-shortener/config"
	"github.com/nelsonandreass/url-shortener/db"
	"github.com/nelsonandreass/url-shortener/helper"
	"github.com/nelsonandreass/url-shortener/models"
)

func Login(c *gin.Context) {
	var user, request models.User

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
	}

	if err := db.DB.Where("user_name = ? ", request.UserName).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
		return
	}

	checkPassword := helper.CheckHashBcrypt(request.Password, user.Password)

	if !checkPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "Invalid Password"})
		return
	}

	expiredTime := time.Now().Add(5 * time.Minute)
	claims := &models.Claims{
		UserName: request.UserName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiredTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.JwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func Register(c *gin.Context) {
	var request models.User
	if err := c.ShouldBindBodyWithJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
	}
	hashedPassword, err := helper.HashBcrypt(request.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	request.Password = hashedPassword
	db.DB.Create(&request)
	c.JSON(http.StatusOK, request)
}
