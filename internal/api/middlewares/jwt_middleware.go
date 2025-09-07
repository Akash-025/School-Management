package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("----JWT----")

		token, err := r.Cookie("Bearer")
		fmt.Println("TOken:", token)
		if err != nil {
			http.Error(w, "Authorization Header missing", http.StatusUnauthorized)
			return
		}

		jwt_secret := os.Getenv("JWT_SECRET")

		parsedToken, err := jwt.Parse(token.Value, func(token *jwt.Token) (any, error) {
			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte(jwt_secret), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				fmt.Println("Token expired") // For env JWT_EXPIRES_IN
				http.Error(w, "Token expired", http.StatusUnauthorized)
				return
			}

			http.Error(w, "Token parsing error", http.StatusInternalServerError)
			return
		}

		if parsedToken.Valid {
			fmt.Println("Valid JWT")
		} else {
			fmt.Println("Invalid JWT")
		}

		if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
			fmt.Println(claims["uid"], claims["exp"], claims["role"])
		} else {
			fmt.Println(err)
		}

		next.ServeHTTP(w, r)
	})
}
