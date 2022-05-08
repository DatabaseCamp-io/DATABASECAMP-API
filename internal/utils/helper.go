package utils

import (
	"fmt"
	"math/rand"
	"net/mail"
	"reflect"

	m "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

func HashAndSalt(pwd string) string {
	bytePassword := []byte(pwd)
	hash, err := bcrypt.GenerateFromPassword(bytePassword, 4)
	if err != nil {
		panic(err.Error())
	}
	return string(hash)
}

func ComparePasswords(hashedPwd string, plainPwd string) bool {
	bytePlainPassword := []byte(plainPwd)
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePlainPassword)
	return err == nil
}

func IsEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func IsSqlDuplicateError(err error) bool {
	sqlError, ok := err.(*m.MySQLError)
	return ok && sqlError.Number == 1062
}

func GetKeyList(value map[string]interface{}) (result []string) {
	for k := range value {
		result = append(result, k)
	}
	return
}

func Shuffle(slice interface{}) {
	rv := reflect.ValueOf(slice)
	swap := reflect.Swapper(slice)
	length := rv.Len()
	for i := length - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		swap(i, j)
	}
}

func ToStrings(strs []interface{}) string {
	result := ""

	if len(strs) == 0 {
		return result
	}

	for _, str := range strs {
		result += fmt.Sprintf("'%v',", str)
	}

	return result[:len(result)-1]
}
