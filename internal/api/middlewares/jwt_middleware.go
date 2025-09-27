package middlewares

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"restapi/pkg/utils"

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

		// JWT from HabibWithGO

		// parsedJwt := strings.Split(token.Value, ".")

		// if len(parsedJwt) != 3{
		// 	http.Error(w, "Authorization Header missing, Unauthorized", http.StatusUnauthorized)
		// 	return
		// }
		// jwtHeader := parsedJwt[0]
		// jwtPayload := parsedJwt[1]
		// jwtSignature := parsedJwt[2]

		// message := jwtHeader + "." + jwtPayload
		// byteArrMsg := []byte(message)
		// byteArrSec := []byte(jwt_secret)

		// h := hmac.New(sha256.New, byteArrSec)
		// h.Write(byteArrMsg)
		// hash := h.Sum(nil)
		// newSignature := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(hash)

		// if newSignature != jwtSignature{
		// 	http.Error(w, "Authorization Header missing, Unauthorized", http.StatusUnauthorized)
		// 	return
		// }
		
		// JWT from Habib_With_GO

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

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if ok {
			fmt.Println(claims["uid"], claims["exp"], claims["role"])
		} else {
			fmt.Println(err)
		}

		ctx := context.WithValue(r.Context(), utils.ContextKey("role"), claims["role"])
		ctx = context.WithValue(ctx, utils.ContextKey("expiresAt"), claims["exp"])
		ctx = context.WithValue(ctx, utils.ContextKey("username"), claims["user"])
		ctx = context.WithValue(ctx, utils.ContextKey("usserId"), claims["uid"])

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
