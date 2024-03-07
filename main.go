package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Product struct {
	ID            int     `json:"product_id"`
	Name          string  `json:"product_name"`
	Price         float32 `json:"price"`
	StockQuantity int     `json:"stock_quantity"`
	Description   string  `json:"description"`
}

func main() {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := fiber.New()

	app.Use(cors.New())

	app.Get("/products", func(c *fiber.Ctx) error {
		rows, err := db.Query("SELECT * FROM products")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var products []Product

		for rows.Next() {
			var product Product

			err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.StockQuantity, &product.Description)
			if err != nil {
				log.Fatal(err)
			}

			products = append(products, product)
		}

		return c.JSON(products)
	})

	err = app.Listen(":8000")
	if err != nil {
		log.Fatal(err)
	}
}
