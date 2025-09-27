package router

import (
	"net/http"
	"restapi/internal/api/handlers"
)

func TeachersRouter(mux *http.ServeMux) {

	mux.HandleFunc("GET /teachers", handlers.GetTeachersHandler)
	mux.HandleFunc("GET /teachers/{id}", handlers.GetOneTeacherHandler)

	mux.HandleFunc("POST /teachers/", handlers.AddTeachersHandler)

	mux.HandleFunc("PATCH /teachers", handlers.PatchTeachersHandler)
	mux.HandleFunc("DELETE /teacher/{id}", handlers.DeleteTeacherHandler)
}
