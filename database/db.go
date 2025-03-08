package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	dsn := "host=localhost user=myuser password=mypass dbname=mydb port=5432 sslmode=disable TimeZone=Asia/Tehran"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Error while connecting to DB:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Error while pinging DB:", err)
	}

	fmt.Println("Successfully connected to DB...")
	DB = db

	createTables()
}

func createTables() {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE TABLE IF NOT EXISTS articles (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		content TEXT NOT NULL,
		author_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("Error while creating tables:", err)
	}
}
