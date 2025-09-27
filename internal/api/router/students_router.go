package router

import (
	"net/http"
	"restapi/internal/api/handlers"
)

func StudentsRouter(mux *http.ServeMux) {

	mux.HandleFunc("GET /students", handlers.GetStudentsHandler)
	mux.HandleFunc("GET /students/{id}", handlers.GetOneStudentHandler)

	mux.HandleFunc("POST /students/", handlers.AddStudentsHandler)

	mux.HandleFunc("PATCH /students", handlers.PatchStudentsHandler)
	mux.HandleFunc("DELETE /student/{id}", handlers.DeleteStudentHandler)

	mux.HandleFunc("GET /teacher/{id}/students", handlers.GetStudentsByTeacherId)

}
