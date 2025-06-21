package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Product struct {
	ID    int     `json:"product_id"`
	Name  string  `json:"product_name"`
	Price float32 `json:"product_price"`
	Stock int     `json:"product_stock"`
	Code  string  `json:"product_code"`
	Sells int     `json:"product_sells"`
}

type ProductsJson struct {
	Products []Product `json:"products"`
}

var db *sql.DB

const cachePath = "./cache/top_sellers.json"
const databasePath = "./database.sql"

func main() {
	start := time.Now()
	insert := false
	if len(os.Args) > 1 {
		insert = os.Args[1] == "insert"
	}
	connectDb(insert)
	defer db.Close()
	products := getCachedTopSellers()
	if products == nil {
		products = getTopSellers()
		saveCache(products)
	}
	fmt.Printf("time elapsed: %vs\n", time.Since(start).Seconds())
	for _, p := range products {
		fmt.Println(p)
	}
}

func connectDb(insert bool) {
	fmt.Println("trying to connect db...")
	database, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	db = database
	err = db.Ping()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	createTable := `
		CREATE TABLE IF NOT EXISTS products (
			id 	INTEGER PRIMARY KEY AUTOINCREMENT,
			name 	TEXT NOT NULL,
			price REAL NOT NULL,
			code  TEXT UNIQUE NOT NULL,
			stock INTEGER NOT NULL,
			sells INTEGER NOT NULL
		);`

	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if !insert {
		return
	}

	populateDatabase()
}

func populateDatabase() {
	insertProduct := func(name string, price float32, code string, stock, sells int) {
		_, err := db.Exec("INSERT INTO products(name, price, code, stock, sells) VALUES(?, ?, ?, ?, ?)", name, price, code, stock, sells)
		if err != nil {
			log.Printf("error: %v", err)
			return
		}
		fmt.Println("Produto criado.")
	}

	for _, p := range CreateProducts() {
		insertProduct(p.Name, p.Price, p.Code, p.Stock, p.Sells)
	}
}

func getCachedTopSellers() []Product {
	fmt.Println("getting products in cache")
	file, err := os.Stat(cachePath)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	if time.Since(file.ModTime()) > 5*time.Minute {
		return nil
	}

	productsBytes, err := os.ReadFile(cachePath)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var cachedProducts ProductsJson
	if err = json.Unmarshal(productsBytes, &cachedProducts); err != nil {
		fmt.Println(err)
		return nil
	}

	return cachedProducts.Products
}

func saveCache(products []Product) {
	fmt.Println("getting products in database")
	productsJson := ProductsJson{
		Products: products,
	}

	jsonData, err := json.Marshal(productsJson)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = os.WriteFile(cachePath, jsonData, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func getTopSellers() []Product {
	rows, err := db.Query("SELECT * FROM products ORDER BY sells LIMIT 10;")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	var products []Product
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Code, &p.Sells, &p.Sells)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		log.Fatalf("error: %v", err)
	}

	return products
}
