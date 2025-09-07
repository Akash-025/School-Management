package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	"golang.org/x/crypto/argon2"
)

func HashPassword(pass string, w http.ResponseWriter) (string, error) {
	password := pass
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		http.Error(w, "Error to generate salt", http.StatusBadRequest)
		return "", nil
	}
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	base64Hash := base64.StdEncoding.EncodeToString(hash)
	base64Salt := base64.StdEncoding.EncodeToString(salt)

	return fmt.Sprintf("%s.%s", base64Salt, base64Hash), nil
}
