package db

import (
	"ecommerce/category/model"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDB() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Database Connection Failed !")
	}
	connstring := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=5432 sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"))
	db, err := gorm.Open(postgres.Open(connstring), &gorm.Config{})
	if err != nil {
		log.Fatalf("Database Connection Failed !")
	}

	Migrate(db)

	return db
}

func Migrate(db *gorm.DB) {

	err := db.AutoMigrate(&model.Category{})
	if err != nil {
		log.Fatalf("migration failed !")
	}

}