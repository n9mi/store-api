package util

import (
	"golang.org/x/crypto/bcrypt"
)

func HashUserPassword(plainPwd string) (string, error) {
	result, err := bcrypt.GenerateFromPassword([]byte(plainPwd), bcrypt.DefaultCost)

	return string(result), err
}

func IsUserPasswordValid(hashPwd string, plainPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPwd), []byte(plainPwd))

	return err == nil
}
