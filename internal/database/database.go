package database

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name  string
	Email string
}

func InitDB() (*gorm.DB, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading env file %v", err)

	}
	// use different db for different environments
	dbName := "development.db"
	if env := os.Getenv("ENV"); env != "" {
		dbName = env + ".db"
	}
	dbPath := filepath.Join(".", dbName)

	// check if database file exists, if not create it and run migrations
	isNewDB := false
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		file, err := os.Create(dbPath)
		if err != nil {
			log.Fatalf("failed to create the database file: %v", err)
		}
		file.Close()
		isNewDB = true
	}

	// connect to the SQLite database
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database")
	}
	// run migrations if it is a new database
	if isNewDB {
		err = runMigrations(db)
		if err != nil {
			log.Fatalf("failed to run migrations: %v", err)
		}
	}
	return db, nil

}
func runMigrations(db *gorm.DB) error {
	// auto migrate the models
	return db.AutoMigrate(
		&User{}, //add other models here
	)
}

// seed data
func SeedDevelopmentData(db *gorm.DB) error {
	users := []*User{
		{Name: "Test User", Email: "test@example.com"},
		{Name: "Another User", Email: "another@example.com"},
	}
	return db.Create(users).Error
}
