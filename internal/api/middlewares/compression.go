package middlewares

import (
	"compress/gzip"
	"fmt"
	"net/http"
	"strings"
)

func CompressionMiddlware(next http.Handler) http.Handler{
	fmt.Println("Compression Middleware")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

        fmt.Println("Compression middlware being returned...")
		
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip"){
			next.ServeHTTP(w, r)

		}

		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()

		w = &gzipResponseWriter{ResponseWriter: w, writer: gz}
		

		next.ServeHTTP(w, r)
		fmt.Println("Compression middleware ends")

	})
}

type gzipResponseWriter struct{
	http.ResponseWriter
	writer *gzip.Writer
}

func(g *gzipResponseWriter) Write(b []byte)(int, error){
	return g.writer.Write(b)
}