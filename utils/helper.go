package utils

import (
	"net/mail"

	m "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type helper struct{}

func NewHelper() helper {
	return helper{}
}

func (h helper) HashAndSalt(pwd string) string {
	bytePassword := []byte(pwd)
	hash, err := bcrypt.GenerateFromPassword(bytePassword, 4)
	if err != nil {
		panic(err.Error())
	}
	return string(hash)
}

func (h helper) ComparePasswords(hashedPwd string, plainPwd string) bool {
	bytePlainPassword := []byte(plainPwd)
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePlainPassword)
	return err == nil
}

func (h helper) IsEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func (h helper) IsSqlDuplicateError(err error) bool {
	sqlError, ok := err.(*m.MySQLError)
	return ok && sqlError.Number == 1062
}

func (h helper) GetKeyList(value map[string]interface{}) (result []string) {
	for k := range value {
		result = append(result, k)
	}
	return
}
