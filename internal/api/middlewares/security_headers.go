package middlewares

import (
	"fmt"
	"net/http"
)


func SecurityHeaders(next http.Handler) http.Handler {
	fmt.Println("Security Header middlware")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Security Header middlware being returned...")

		w.Header().Set("X-DNS-Prefetch-Control", "off")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("X-Powered-By", "Django")

		next.ServeHTTP(w, r)
		fmt.Println("Security Header middlware ends")
	})
}


//  ===> Basic Middlewares skeleton <==
// func SecurityHeaders(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		next.ServeHTTP(w, r)
// 	})
// }
