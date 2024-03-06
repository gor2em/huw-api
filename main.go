package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var (
	DB                           *sql.DB
	PRODUCTION_APP_ENV           string = "production"
	DEVELOPMENT_ENVIRONMENT_FILE string = ".env.dev"
)

type Product struct {
	ID            int     `json:"product_id"`
	Name          string  `json:"product_name"`
	Price         float32 `json:"price"`
	StockQuantity int     `json:"stock_quantity"`
	Description   string  `json:"description"`
}

func init() {
	initEnv()

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dbSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	var errDB error
	DB, errDB = sql.Open("mysql", dbSource)
	if errDB != nil {
		log.Fatal(errDB)
	}

	errDB = DB.Ping()
	if errDB != nil {
		log.Fatal(errDB)
	}

}

func initEnv() {
	env := os.Getenv("APP_ENV")

	if env != PRODUCTION_APP_ENV {

		if err := godotenv.Load(DEVELOPMENT_ENVIRONMENT_FILE); err != nil {
			fmt.Printf("failed to load environment variables: %v", err)
			return
		}

		fmt.Println("=== LOADED ===", DEVELOPMENT_ENVIRONMENT_FILE)
	} else {
		fmt.Println("=== LOADED === ", PRODUCTION_APP_ENV)
	}
}

func main() {
	http.HandleFunc("/products", getProducts)

	server := &http.Server{
		Addr:         ":" + os.Getenv("APP_PORT"),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Server listening on :%s", os.Getenv("APP_PORT"))
	log.Fatal(server.ListenAndServe())
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query("SELECT * FROM products")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.StockQuantity, &product.Description)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
