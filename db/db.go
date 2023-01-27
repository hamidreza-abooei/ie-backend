package db

import (
	"fmt"
	"strings"

	"github.com/hamidreza-abooei/ie-project/model"
	"github.com/jinzhu/gorm"
)

// Personally I prefre using SQL Server instead of SQLite but it's simpler so lets use it for now

// Database setup
func Setup(dbName string) *gorm.DB {
	// Create Database
	db := CreateDB(dbName)
	migrate(db)

	return db
}

// Create Database
func CreateDB(dbName string) *gorm.DB {
	// Fix Name
	if !strings.HasSuffix(dbName, ".db") {
		dbName = dbName + ".db"
	}
	// Open DB (Or create if not exists)
	db, err := gorm.Open("sqlite3", "./"+dbName)
	// Error handling
	if err != nil {
		fmt.Println("Error in opening Sqlite:", err)
		return nil
	}
	return db
}

func migrate(db *gorm.DB) {

	db.AutoMigrate(&model.User{}, &model.Url{}, &model.Request{})
}
