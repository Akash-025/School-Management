package handlers

import (
	"fmt"
	"net/http"
)

func ExecsHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:

		// fmt.Println("Query: ", r.URL.Query())
		// fmt.Println("name: ", r.URL.Query().Get("name"))

		// err := r.ParseForm()
		// if err != nil {
		// 	return
		// }

		// fmt.Println("Form from Post method: ", r.Form)
		w.Write([]byte("Hello post method on execs route"))

	case http.MethodGet:
		w.Write([]byte("Hello get method on execs route"))
		//fmt.Println("Hello get method on execs route")
	case http.MethodPut:
		w.Write([]byte("Hello put method on execs route"))
		fmt.Println("Hello put method on execs route")
	case http.MethodPatch:
		w.Write([]byte("Hello patch method on execs route"))
		fmt.Println("Hello patch method on execs route")
	case http.MethodDelete:
		w.Write([]byte("Hello delete method on execs route"))
		fmt.Println("Hello delete method on execs route")

	}
}