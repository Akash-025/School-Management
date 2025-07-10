package main

import (
	
	"fmt"
	"net/http"
	"strings"
)

type user struct{
	Name string `json:"name"`
	Age string `json:"age"`
	City string `jsong:"city"`
}


func rootHandler( w http.ResponseWriter, r *http.Request)  {
		w.Write([]byte("Hello Root route"))
		fmt.Println("Hello root route")
		fmt.Println(r.Method)
	}

func studentHandler( w http.ResponseWriter, r *http.Request)  {
		
		switch r.Method {
		case http.MethodPost :
			w.Write([]byte("Hello get method on students route"))
		    fmt.Println("Hello get method on students route")
		case http.MethodGet :
			w.Write([]byte("Hello post method on students route"))
		    fmt.Println("Hello post method on students route")
		case http.MethodPut :
			w.Write([]byte("Hello put method on students route"))
		    fmt.Println("Hello put method on students route")
		case http.MethodPatch :
			w.Write([]byte("Hello patch method on students route"))
		    fmt.Println("Hello patch method on students route")
		case http.MethodDelete :
			w.Write([]byte("Hello delete method on students route"))
		    fmt.Println("Hello delete method on students route")

		}
	}

func teachersHandler( w http.ResponseWriter, r *http.Request)  {
		
		switch r.Method {
		case http.MethodGet :

			path := strings.TrimPrefix(r.URL.Path, "/teachers/")
			pathID := strings.TrimSuffix(path, "/")
			fmt.Println("ID is: ", pathID)

			queryParams := r.URL.Query()
			sortBy := queryParams.Get("sortBy")
			key := queryParams.Get("key")
			sortorder := queryParams.Get("sortorder")

			fmt.Printf("SortBy: %v, Key: %v, SortOrder: %v\n", sortBy, key, sortorder)
			

			w.Write([]byte("Hello get method on teachers route"))
		    fmt.Println("Hello get method on teachers route")
		case http.MethodPost :
			w.Write([]byte("Hello post method on teachers route"))
		    fmt.Println("Hello post method on teachers route")
		case http.MethodPut :
			w.Write([]byte("Hello put method on teachers route"))
		    fmt.Println("Hello put method on teachers route")
		case http.MethodPatch :
			w.Write([]byte("Hello patch method on teachers route"))
		    fmt.Println("Hello patch method on teachers route")
		case http.MethodDelete :
			w.Write([]byte("Hello delete method on teachers route"))
		    fmt.Println("Hello delete method on teachers route")

		}
	}	

func execsHandler( w http.ResponseWriter, r *http.Request)  {
		
		switch r.Method {
		case http.MethodPost :
			w.Write([]byte("Hello post method on execs route"))
		    fmt.Println("Hello post method on execs route")
		case http.MethodGet :
			w.Write([]byte("Hello get method on execs route"))
		    fmt.Println("Hello get method on execs route")
		case http.MethodPut :
			w.Write([]byte("Hello put method on execs route"))
		    fmt.Println("Hello put method on execs route")
		case http.MethodPatch :
			w.Write([]byte("Hello patch method on execs route"))
		    fmt.Println("Hello patch method on execs route")
		case http.MethodDelete :
			w.Write([]byte("Hello delete method on execs route"))
		    fmt.Println("Hello delete method on execs route")

		}
	}
	

func main() {
	port := ":3000"

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/students/", studentHandler)
	http.HandleFunc("/teachers/", teachersHandler)
	http.HandleFunc("/execs/", execsHandler)

	fmt.Println("Server is running on port",port)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		
	}
}