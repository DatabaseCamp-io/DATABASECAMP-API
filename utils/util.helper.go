package utils

import (
	"math/rand"
	"reflect"
	"regexp"

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
	emailRegExp := regexp.MustCompile(`/^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/`)
	return emailRegExp.MatchString(email)
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

func (h helper) Shuffle(slice interface{}) {
	rv := reflect.ValueOf(slice)
	swap := reflect.Swapper(slice)
	length := rv.Len()
	for i := length - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		swap(i, j)
	}
}
