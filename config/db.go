package config

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var DB *sql.DB

func init() {
	// โหลด Environment Variables จากไฟล์ .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

// getEnv retrieves environment variables or returns a default value if not found
func getEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("Environment variable %s not set", key)
	}
	return value
}

func Connect() {
	var err error
	databaseURL := getEnv("DATABASE_URL")
	print(databaseURL)
	DB, err = sql.Open("mysql", databaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Connected to database")
}
