package db

import (
	"fmt"
	"log"
	"os"

	"github.com/nelsonandreass/url-shortener/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, name)
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed connect database", err)
	}
	database.AutoMigrate(&models.URL{})
	database.AutoMigrate(&models.User{})

	err = database.AutoMigrate(
		&models.URL{},
		&models.User{},
	)

	if err != nil {
		log.Fatal("failed to migrate models: ", err)
	}

	DB = database
}
