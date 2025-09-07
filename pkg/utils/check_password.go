package utils

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"golang.org/x/crypto/argon2"
)

func PasswordCheck(passFromDb string, w http.ResponseWriter, inputPassword string) (string, bool) {
	passwordFromDb := passFromDb
	fmt.Println("Password: ", passwordFromDb)
	parts := strings.Split(passwordFromDb, ".")
	if len(parts) != 2 {
		http.Error(w, "Password split error", http.StatusInternalServerError)
		return "", true
	}
	salt := parts[0]
	hash := parts[1]

	saltByte, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		http.Error(w, "SaltByte error", http.StatusInternalServerError)
		return "", true
	}
	hashByte, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		http.Error(w, "HashByte error", http.StatusInternalServerError)
		return "", true
	}
	hashFromDb := argon2.IDKey([]byte(inputPassword), saltByte, 1, 64*1024, 4, 32)

	var msg string
	if len(hashFromDb) != len(hashByte) {
		msg = "Incorrect password"
		return "", true
	}

	if subtle.ConstantTimeCompare(hashFromDb, hashByte) == 1 {
		msg = "Successfully Logged in"

	} else {
		msg = "Incorrect password"
		return "", true

	}
	return msg, false
}
