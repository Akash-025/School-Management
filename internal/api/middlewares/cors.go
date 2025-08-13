package middlewares

import (
	"fmt"
	"net/http"
)

var allowedOrigins = []string{
	"https://myfrontend.com",
	"https://localhost:3000",
}

func Cors(next http.Handler) http.Handler{
	fmt.Println("CORS middlware")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("CORS middlware being returned...")
		origin := r.Header.Get("origin")

		if isAllowedOrigin(origin){
			w.Header().Set("Access-Control-Allow", origin)
		}else{
			http.Error(w, "Not allowed by CORS", http.StatusForbidden)
			return
		}

		if r.Method == http.MethodOptions{
			return
		}

		next.ServeHTTP(w, r)
		fmt.Println("CORS middleware ends")
	})
}

func isAllowedOrigin(origin string) bool{
	for _,val := range allowedOrigins{
		if val == origin{
			return true
		}
	}
	return false
}