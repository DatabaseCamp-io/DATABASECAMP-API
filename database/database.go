package database

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type db struct {
	db *gorm.DB
}

type IDatabase interface {
	GetDB() *gorm.DB
	CloseDB() error
	OpenConnection() error
}

var instantiated *db = nil

func New() *db {
	if instantiated == nil {
		instantiated = new(db)
	}
	return instantiated
}

func (database db) getDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
}

func (database db) OpenConnection() error {
	var err error
	dsn := database.getDSN()
	sql := mysql.Open(dsn)
	database.db, err = gorm.Open(sql)
	return err
}

func (database db) GetDB() *gorm.DB {
	return database.db
}

func (database db) CloseDB() error {
	sql, err := database.db.DB()
	if err != nil {
		return err
	}
	err = sql.Close()
	return err
}
