package database

// database.go
/**
 * 	This file used to create database transaction
 */

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/**
 * This class manage database transaction
 */
type transaction struct {
	db *gorm.DB
}

/**
 * Constructor creates a new transaction instance
 *
 * @return 	instance of transaction
 */
func NewTransaction() *transaction {
	dsn := getDSN()
	sql := mysql.Open(dsn)
	db, _ := gorm.Open(sql)
	return &transaction{db: db}
}

/**
 * Get the database for data manipulation
 *
 * @return the database that can be used to manipulate data in the database
 */
func (t *transaction) GetDB() *gorm.DB {
	return t.db
}

/**
 * Begin database transaction
 */
func (t *transaction) Begin() {
	t.db = t.db.Begin()
}

/**
 * Commit database transaction
 */
func (t *transaction) Commit() {
	t.db = t.db.Commit()
}

/**
 * Rollback database transaction
 */
func (t *transaction) Rollback() {
	t.db = t.db.Rollback()
}

/**
 * Close database transaction
 *
 * @return the error of closing the database transaction
 */
func (t *transaction) Close() error {
	sql, err := t.db.DB()
	if err != nil {
		return err
	}
	err = sql.Close()
	return err
}
