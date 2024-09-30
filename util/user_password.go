package util

import "golang.org/x/crypto/bcrypt"

func NewUserPassword(plainPwd string) (string, error) {
	result, err := bcrypt.GenerateFromPassword([]byte(plainPwd), bcrypt.DefaultCost)

	return string(result), err
}
