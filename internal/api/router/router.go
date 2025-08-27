package router

import (
	"net/http"
	"restapi/internal/api/handlers"
)

func Router() *http.ServeMux {

	mux := http.NewServeMux()
// FOR TEACHERS
	mux.HandleFunc("GET /teachers", handlers.GetTeachersHandler)
	mux.HandleFunc("GET /teachers/{id}", handlers.GetOneTeacherHandler)

	mux.HandleFunc("POST /teachers/", handlers.AddTeachersHandler)

	mux.HandleFunc("PATCH /teachers", handlers.PatchTeachersHandler)
	mux.HandleFunc("DELETE /teacher/{id}", handlers.DeleteTeacherHandler)


// FOR STUDENTS
	mux.HandleFunc("GET /students", handlers.GetStudentsHandler)
	mux.HandleFunc("GET /students/{id}", handlers.GetOneStudentHandler)

	mux.HandleFunc("POST /students/", handlers.AddStudentsHandler)

	mux.HandleFunc("PATCH /students", handlers.PatchStudentsHandler)
	mux.HandleFunc("DELETE /student/{id}", handlers.DeleteStudentHandler)


	return mux

}
