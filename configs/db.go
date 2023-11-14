package configs

import (
	"context"

	"fmt"
	"log"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "58200Mm#"
	dbname   = "ecommerce2"
)

func ConnectToDB() *gorm.DB {
	connstring := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%d sslmode=disable", user, dbname, password, host, port)
	db, err := gorm.Open(postgres.Open(connstring), &gorm.Config{})
	if err != nil {
		log.Fatalf("Database Connection Failed !")
	}

	Migrate(db)

	return db
}

func Migrate(db *gorm.DB) {

	err := db.AutoMigrate(&User{})
	if err != nil {
		log.Fatalf("migration failed !")
	}
	err = db.AutoMigrate(&Category{})
	if err != nil {
		log.Fatalf("migration failed !")
	}
	err = db.AutoMigrate(&Product{})
	if err != nil {
		log.Fatalf("migration failed !")
	}
	err = db.AutoMigrate(&Order{})
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
