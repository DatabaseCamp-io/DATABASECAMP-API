package database

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type database struct {
	db *gorm.DB
}

type IDatabase interface {
	GetDB() *gorm.DB
	CloseDB() error
	OpenConnection() error
}

var instantiated *database = nil

// Create database instance
func New() *database {
	if instantiated == nil {
		instantiated = new(database)
	}
	return instantiated
}

// Get DSN
func getDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
}

// Connect to database
func (db *database) OpenConnection() error {
	var err error
	dsn := getDSN()
	sql := mysql.Open(dsn)
	db.db, err = gorm.Open(sql)
	return err
}

// Get DB variable
func (db *database) GetDB() *gorm.DB {
	return db.db
}

// Close DB
func (db *database) CloseDB() error {
	sql, err := db.db.DB()
	if err != nil {
		return err
	}
	err = sql.Close()
	return err
}
