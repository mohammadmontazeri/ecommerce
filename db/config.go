package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "58200Mm#"
	dbname   = "ecommerce"
)

func ConnectToDb() *sql.DB {

	connstring := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%d sslmode=disable", user, dbname, password, host, port)
	db, err := sql.Open("postgres", connstring)

	if err != nil {
		log.Fatalf("Database Connection Failed !")
	} else {
		fmt.Println("You are connect to DB !")

	}

	return db
}

func MigrateTables(c *gin.Context) {

	db := ConnectToDb()

	userTable := `CREATE TABLE users (
		id SERIAL PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		password TEXT,
		email TEXT UNIQUE NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	  )`
	categoryTable := `CREATE TABLE categories (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) UNIQUE NOT NULL,
		parent_id INT,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		CONSTRAINT fk_parent_id
      	FOREIGN KEY(parent_id)
      	REFERENCES categories(id)
		ON UPDATE CASCADE 
		ON DELETE CASCADE 
	  )`

	productTable := `CREATE TABLE products (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) UNIQUE NOT NULL,
		code VARCHAR(100) UNIQUE NOT NULL ,
		price DECIMAL,
		picture TEXT ,
		detail TEXT ,
		category_id INT,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		CONSTRAINT fk_category_id
      	FOREIGN KEY(category_id)
      	REFERENCES categories(id)
	  )`

	orderTable := `CREATE TABLE orders (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL,
		product_id INT NOT NULL,
		code VARCHAR(100) NOT NULL ,
		price DECIMAL NOT NULL,
		status INT NUT NULL,
		CONSTRAINT fk_user_id
      	FOREIGN KEY(user_id)
      	REFERENCES users(id),
		CONSTRAINT fk_product_id
      	FOREIGN KEY(product_id)
      	REFERENCES products(id)
		`
	// new table added in tables map

	tables := map[string]string{
		"userTable":     userTable,
		"categoryTable": categoryTable,
		"productTable":  productTable,
		"orderTable":    orderTable,
	}

	for _, value := range tables {
		db.Exec(value)
	}

}
