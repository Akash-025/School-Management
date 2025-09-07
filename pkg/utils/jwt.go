package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func SignToken(userId int, username, role string) (string, error) {

	jwt_secret := os.Getenv("JWT_SECRET")
	jwt_expires_in := os.Getenv("JWT_EXPIRES_IN")

	claims := jwt.MapClaims{
		"uId":   userId,
		"uname": username,
		"role":  role,
	}

	if jwt_expires_in != "" {
		duration, err := time.ParseDuration(jwt_expires_in)
		if err != nil {
			return "", nil
		}
		claims["exp"] = jwt.NewNumericDate(time.Now().Add(duration))
	} else {
		claims["exp"] = jwt.NewNumericDate(time.Now().Add(15 * time.Minute))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwt_secret))
	if err != nil {
		return "", nil
	}
	return signedToken, nil
}
