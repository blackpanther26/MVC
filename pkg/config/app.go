package config

import (
	"log"
	"os"
	"github.com/blackpanther26/mvc/pkg/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var	DB *gorm.DB

func ConnectToDb()  {
	dsn := os.Getenv("DB")
	var err error
	DB ,err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	log.Println("Successfully connected to the database")
}

func SyncDatabase() {
	err := DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}
	log.Println("Database migration completed")
}

func GetDB() *gorm.DB {
	return DB
}