package database

// database.go
/**
 * 	This file used to create connection to the RDBMS
 */

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/**
 * This class manage database connection
 */
type database struct {
	db *gorm.DB
}

// Instance of database class for singleton pattern
var instantiated *database = nil

/**
 * Constructor creates a new database instance or geting a database instance
 *
 * @return 	instance of database
 */
func New() *database {
	if instantiated == nil {
		instantiated = new(database)
	}
	return instantiated
}

/**
 * Get a name of the database from the environment
 *
 * @return name of the database
 */
func getDBName() string {
	if os.Getenv("MODE") == "develop" {
		return os.Getenv("DB_NAME_DEVELOP")
	} else {
		return os.Getenv("DB_NAME_PRODUCTION")
	}
}

/**
 * Get a DSN of the database from the environment
 *
 * @return DSN of the database
 */
func getDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		getDBName(),
	)
}

/**
 * Get the database for data manipulation
 *
 * @return the database that can be used to manipulate data in the database
 */
func (db *database) GetDB() *gorm.DB {
	return db.db
}

/**
 * Open the database connection
 *
 * @return the error of opening the database
 */
func (db *database) OpenConnection() error {
	var err error
	dsn := getDSN()
	sql := mysql.Open(dsn)
	db.db, err = gorm.Open(sql)
	return err
}

/**
 * Close the database connection
 *
 * @return the error of closing the database
 */
func (db *database) CloseDB() error {
	sql, err := db.db.DB()
	if err != nil {
		return err
	}
	err = sql.Close()
	return err
}
