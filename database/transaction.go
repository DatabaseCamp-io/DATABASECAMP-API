package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type transaction struct {
	db *gorm.DB
}

type ITransaction interface {
	GetDB() *gorm.DB
	Begin()
	Commit()
	Rollback()
	Close() error
}

// Create transaction instance
func NewTransaction() *transaction {
	dsn := getDSN()
	sql := mysql.Open(dsn)
	db, _ := gorm.Open(sql)
	return &transaction{db: db}
}

// Get database
func (t *transaction) GetDB() *gorm.DB {
	return t.db
}

// Start database
func (t *transaction) Begin() {
	t.db = t.db.Begin()
}

// Commit database
func (t *transaction) Commit() {
	t.db = t.db.Commit()
}

// Rollback database
func (t *transaction) Rollback() {
	t.db = t.db.Rollback()
}

// Close database
func (t *transaction) Close() error {
	sql, err := t.db.DB()
	if err != nil {
		return err
	}
	err = sql.Close()
	return err
}
