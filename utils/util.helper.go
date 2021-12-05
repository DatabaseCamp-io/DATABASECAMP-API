package utils

// util.helper.go
/**
 * 	This file is a part of utilities, used collect normal functions of the application
 */

import (
	"math/rand"
	"net/mail"
	"reflect"

	m "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

/**
 * 	This class normal functions of the application
 */
type helper struct{}

/**
 * Constructor creates a new helper instance
 *
 * @return 	instance of helper
 */
func NewHelper() helper {
	return helper{}
}

/**
 * Hash password with salt
 *
 * @param	paw  	Raw password to hash
 *
 * @return hashed password
 */
func (h helper) HashAndSalt(pwd string) string {
	bytePassword := []byte(pwd)
	hash, err := bcrypt.GenerateFromPassword(bytePassword, 4)
	if err != nil {
		panic(err.Error())
	}
	return string(hash)
}

/**
 * Compare password with hashed password
 *
 * @param	hashedPwd  	Hashed password
 * @param	plainPwd  	Raw password to compare
 *
 * @return true if equal, false if not
 */
func (h helper) ComparePasswords(hashedPwd string, plainPwd string) bool {
	bytePlainPassword := []byte(plainPwd)
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePlainPassword)
	return err == nil
}

/**
 * Check email address format
 *
 * @param	Email  Email to validate
 *
 * @return true if email valid, false if not
 */
func (h helper) IsEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

/**
 * Check SQL duplicate error
 *
 * @param	err  Error to check
 *
 * @return true if error is sql duplicate value, false if not
 */
func (h helper) IsSqlDuplicateError(err error) bool {
	sqlError, ok := err.(*m.MySQLError)
	return ok && sqlError.Number == 1062
}

/**
 * Get list of the key of map
 *
 * @param	value  Map to get key
 *
 * @return list of the key of map
 */
func (h helper) GetKeyList(value map[string]interface{}) (result []string) {
	for k := range value {
		result = append(result, k)
	}
	return
}

/**
 * Shuffle member in list
 *
 * @param	slice  List to shuffle
 */
func (h helper) Shuffle(slice interface{}) {
	rv := reflect.ValueOf(slice)
	swap := reflect.Swapper(slice)
	length := rv.Len()
	for i := length - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		swap(i, j)
	}
}
