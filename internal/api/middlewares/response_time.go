package middlewares

import (
	"fmt"
	"net/http"
	"time"

)

func ResponseTimeMiddlware(next http.Handler) http.Handler{
	fmt.Println("Response Time middlware")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Response Time middlware being returned...")

		start := time.Now()
		wrappedWriter := &resposeWriter{ResponseWriter: w, status: http.StatusOK}

		time.Sleep(200*time.Microsecond)
		//Calculate the duration
		duration := time.Since(start)
		w.Header().Set("X-Response-Time", duration.String())
		
		next.ServeHTTP(w, r)
		// Loggin
		fmt.Printf("Method: %s, URL: %s, Status:%d, Duration: %s\n", r.Method, r.URL, wrappedWriter.status, duration.String())
		fmt.Println("Response Time middlware ends")

	})

}

type resposeWriter struct{
	http.ResponseWriter
	status int
}

func (rw *resposeWriter) WriteHeader(code int){
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}