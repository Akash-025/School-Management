package handlers

import (
	"fmt"
	"net/http"
)

func StudentHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		w.Write([]byte("Hello get method on students route"))
		//fmt.Println("Hello get method on students route")
	case http.MethodGet:
		w.Write([]byte("Hello post method on students route"))
		//fmt.Println("Hello post method on students route")
	case http.MethodPut:
		w.Write([]byte("Hello put method on students route"))
		fmt.Println("Hello put method on students route")
	case http.MethodPatch:
		w.Write([]byte("Hello patch method on students route"))
		fmt.Println("Hello patch method on students route")
	case http.MethodDelete:
		w.Write([]byte("Hello delete method on students route"))
		fmt.Println("Hello delete method on students route")

	}
}
