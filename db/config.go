package db

import (
	"ecommerce/category"
	"ecommerce/order"
	"ecommerce/product"
	userPackage "ecommerce/user"
	"fmt"
	"log"

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

	db.AutoMigrate(&userPackage.User{})
	db.AutoMigrate(&category.Category{})
	db.AutoMigrate(&product.Product{})
	db.AutoMigrate(&order.Order{})

}
