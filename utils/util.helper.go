package utils

import (
	"math/rand"
	"net/mail"
	"reflect"

	m "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type helper struct{}

func NewHelper() helper {
	return helper{}
}

// Use to Generate form password
func (h helper) HashAndSalt(pwd string) string {
	bytePassword := []byte(pwd)
	hash, err := bcrypt.GenerateFromPassword(bytePassword, 4)
	if err != nil {
		panic(err.Error())
	}
	return string(hash)
}

// Use to compare password and form
func (h helper) ComparePasswords(hashedPwd string, plainPwd string) bool {
	bytePlainPassword := []byte(plainPwd)
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePlainPassword)
	return err == nil
}

// Usto check validity of email
func (h helper) IsEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// Check duplicate sql
func (h helper) IsSqlDuplicateError(err error) bool {
	sqlError, ok := err.(*m.MySQLError)
	return ok && sqlError.Number == 1062
}

// Use to get list of key
func (h helper) GetKeyList(value map[string]interface{}) (result []string) {
	for k := range value {
		result = append(result, k)
	}
	return
}

// Use to shuffle varieble
func (h helper) Shuffle(slice interface{}) {
	rv := reflect.ValueOf(slice)
	swap := reflect.Swapper(slice)
	length := rv.Len()
	for i := length - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		swap(i, j)
	}
}
