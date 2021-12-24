package database

// interface.go
/**
 * 	This file used to be a interface of database
 */

import "gorm.io/gorm"

/**
 * 	 Interface to show function in database that others can use
 */
type IDatabase interface {

	/**
	 * Get the database for data manipulation
	 *
	 * @return the database that can be used to manipulate data in the database
	 */
	GetDB() *gorm.DB

	/**
	 * Open the database connection
	 *
	 * @return the error of opening the database
	 */
	OpenConnection() error

	/**
	 * Close the database connection
	 *
	 * @return the error of closing the database
	 */
	CloseDB() error
}

/**
 * 	 Interface to show function in transaction that others can use
 */
type ITransaction interface {

	/**
	 * Get the database for data manipulation
	 *
	 * @return the database that can be used to manipulate data in the database
	 */
	GetDB() *gorm.DB

	/**
	 * Begin database transaction
	 */
	Begin()

	/**
	 * Commit database transaction
	 */
	Commit()

	/**
	 * Rollback database transaction
	 */
	Rollback()

	/**
	 * Close database transaction
	 *
	 * @return the error of closing the database transaction
	 */
	Close() error
}
