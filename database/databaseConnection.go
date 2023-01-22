package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func DbInstance() *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbProtocol := os.Getenv("DB_PROTOCOL")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dbConnectionFormat := "%s:%s@%s(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	dbSourceName := fmt.Sprintf(dbConnectionFormat, dbUsername, dbPassword, dbProtocol, dbHost, dbPort, dbName)

	db, err := gorm.Open(mysql.Open(dbSourceName), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if db.WithContext(ctx).Error != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MYSQL")
	return db
}

var DB *gorm.DB = DbInstance()

func ModelForDbOperations(db *gorm.DB, model any) *gorm.DB {
	err := db.AutoMigrate(&model)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Model(&model).Error; err != nil {
		log.Fatal(err)
	}
	return db
}
