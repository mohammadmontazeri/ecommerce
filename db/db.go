package db

import (
	"context"
	"ecommerce/models"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
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

	err := db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("migration failed !")
	}
	err = db.AutoMigrate(&models.Category{})
	if err != nil {
		log.Fatalf("migration failed !")
	}
	err = db.AutoMigrate(&models.Product{})
	if err != nil {
		log.Fatalf("migration failed !")
	}
	err = db.AutoMigrate(&models.Order{})
	if err != nil {
		log.Fatalf("migration failed !")
	}

}

func ConnectToRedis() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		return nil, err
	} else {
		return rdb, err
	}

}

func ConnectToRedisForCache() (*cache.Cache, error) {

	rdb, err := ConnectToRedis()
	if err != nil {
		return nil, err
	}
	myCache := cache.New(&cache.Options{
		Redis: rdb,
	})

	return myCache, nil
}
