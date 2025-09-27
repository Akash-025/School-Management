package router

import (
	"net/http"
	"restapi/internal/api/handlers"
)

func ExecsRouter(mux *http.ServeMux) {

	mux.HandleFunc("GET /execs", handlers.GetExecsHandler)
	mux.HandleFunc("GET /execs/{id}", handlers.GetOneExecHandler)
	mux.HandleFunc("POST /execs", handlers.AddExecsHandler)
	mux.HandleFunc("PATCH /execs", handlers.PatchExecsHandler)
	mux.HandleFunc("DELETE /execs/{id}", handlers.DeleteExecHandler)

	mux.HandleFunc("POST /execs/{id}/updatepassword", handlers.UpdatePasswordExecHandler)

	mux.HandleFunc("POST /execs/login", handlers.LoginExecsHandler)
	mux.HandleFunc("POST /execs/logout", handlers.LogoutExecsHandler)
	mux.HandleFunc("POST /execs/forgotpassword", handlers.PatchStudentsHandler)
	mux.HandleFunc("POST /execs/resetpassword/reset/{resetcode}", handlers.PatchStudentsHandler)
}
