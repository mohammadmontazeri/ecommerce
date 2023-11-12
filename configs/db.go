package configs

import (
	"ecommerce/internal/category"
	"ecommerce/internal/order"
	"ecommerce/internal/product"

	userPackage "ecommerce/internal/user"
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

	// Migrate(db)

	return db
}

func Migrate(db *gorm.DB) {

	db.AutoMigrate(&userPackage.User{})
	db.AutoMigrate(&category.Category{})
	db.AutoMigrate(&product.Product{})
	db.AutoMigrate(&order.Order{})

}

func ConnectToRedisForCache() (*redis.Client, *cache.Cache) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	mycache := cache.New(&cache.Options{
		Redis: rdb,
	})

	return rdb, mycache
}
