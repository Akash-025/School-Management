package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"restapi/internal/api/middlewares"
	"restapi/internal/api/router"
	"restapi/internal/repositories/sqlconnect"
	"restapi/pkg/utils"
	"time"
)

type user struct {
	Name string `json:"name"`
	Age  string `json:"age"`
	City string `jsong:"city"`
}

func main() {

	_, err := sqlconnect.ConnectDb()
	if err != nil {
		fmt.Println("Error-----")
		return
	}

	port := os.Getenv("API_PORT") //From environment variable

	cert := "cert.pem"
	key := "key.pem"

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	rl := middlewares.NewRateLimiter(5, time.Minute)

	hppOptions := middlewares.HPPoptions{
		CheckQuery:              true,
		CheckBody:               true,
		CheckBodyOnlyForContent: "application/x-www-form-urlencoded",
		WhiteList:               []string{"sortOrder", "sortBy", "name"},
	}
	//router := router.Router()
	mux := router.Router()
	//secureMux := middlewares.Cors(rl.Middleware(middlewares.ResponseTimeMiddlware(middlewares.SecurityHeaders(middlewares.CompressionMiddlware(middlewares.Hpp(hppOptions)(mux))))))
	secureMux := utils.ApplyMiddlewares(mux, middlewares.Hpp(hppOptions), middlewares.CompressionMiddlware, middlewares.SecurityHeaders, middlewares.ResponseTimeMiddlware, rl.Middleware, middlewares.Cors)
	//secureMux = middlewares.JWTMiddleware(secureMux)
	jwtMiddleware := middlewares.MiddlewaresExcludePaths(middlewares.JWTMiddleware, "/execs/login")
	secureMux = jwtMiddleware(secureMux)
	//Create custom server
	server := &http.Server{
		Addr:      port,
		Handler:   secureMux,
		TLSConfig: tlsConfig,
	}

	fmt.Println("Server is running on port", port)

	err = server.ListenAndServeTLS(cert, key)
	if err != nil {

	}
}

// Middleware is a function that wraps an http.Handler with additional functionality
