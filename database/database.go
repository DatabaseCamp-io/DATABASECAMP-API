package database

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DB struct {
	db *gorm.DB
}

type IDatabase interface {
	GetDB() *gorm.DB
	CloseDB(db *gorm.DB) error
	OpenConnection() (*gorm.DB, error)
}

var instantiated *DB = nil

func New() *DB {
	if instantiated == nil {
		instantiated = new(DB)
	}
	return instantiated
}

func (database DB) getDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
}

func (database DB) OpenConnection() error {
	var err error
	dsn := database.getDSN()
	sql := mysql.Open(dsn)
	database.db, err = gorm.Open(sql)
	return err
}

func (database DB) Get() *gorm.DB {
	return database.db
}

func (database DB) Close() error {
	sql, err := database.db.DB()
	if err != nil {
		return err
	}
	err = sql.Close()
	return err
}
