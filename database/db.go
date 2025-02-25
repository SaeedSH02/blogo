package database

import (
	"blogo/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := "host=localhost user=myuser password=mypass dbname=mydb port=5432 sslmode=disable TimeZone=Asia/Tehran"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error while connecting to DB....")
	}
	fmt.Println("successfully connected to DB...")
	DB = db

	db.AutoMigrate(&models.User{}, &models.Article{})
}