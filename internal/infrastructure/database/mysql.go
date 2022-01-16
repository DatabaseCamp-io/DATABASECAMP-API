package database

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MysqlDB interface {
	GetDB() *gorm.DB
	OpenConnection() error
	CloseConnection() error
}

type mysqlDB struct {
	db *gorm.DB
}

var instantiated *mysqlDB = nil

func GetMySqlDBInstance() *mysqlDB {
	if instantiated == nil {
		instantiated = new(mysqlDB)
	}
	return instantiated
}

func getDBName() string {
	if os.Getenv("MODE") == "develop" {
		return os.Getenv("DB_NAME_DEVELOP")
	} else {
		return os.Getenv("DB_NAME_PRODUCTION")
	}
}

func getDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		getDBName(),
	)
}

func (db *mysqlDB) GetDB() *gorm.DB {
	return db.db
}

func (db *mysqlDB) OpenConnection() error {
	var err error
	dsn := getDSN()
	sql := mysql.Open(dsn)
	db.db, err = gorm.Open(sql)
	return err
}

func (db *mysqlDB) CloseConnection() error {
	sql, err := db.db.DB()
	if err != nil {
		return err
	}
	err = sql.Close()
	return err
}
