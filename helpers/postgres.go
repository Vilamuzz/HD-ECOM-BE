package helpers

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	uri := "user=" + dbUser + " password=" + dbPassword + " host=" + dbHost + " port=" + dbPort + " dbname=" + dbName + " sslmode=disable"

	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to the database successfully!")
	return db
}

func MigrateDB(db *gorm.DB, models ...interface{}) {
	// Skip migration since we're using existing database
	fmt.Println("Using existing database schema")
}