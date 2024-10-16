package util

import "math/rand"

var letters []byte = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890")
var numLetters = len(letters)

func GenerateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(numLetters)]
	}

	return string(b)
}
