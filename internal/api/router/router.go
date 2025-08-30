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

	mux.HandleFunc("GET /teacher/{id}/students", handlers.GetStudentsByTeacherId)

// For EXECS
    mux.HandleFunc("GET /execs", handlers.GetExecsHandler)
	mux.HandleFunc("GET /execs/{id}", handlers.GetOneExecHandler)
	mux.HandleFunc("POST /execs", handlers.AddExecsHandler)
	mux.HandleFunc("PATCH /execs", handlers.PatchExecsHandler)
	mux.HandleFunc("DELETE /execs/{id}", handlers.DeleteExecHandler)

	mux.HandleFunc("POST /execs/{id}/updatepassword", handlers.DeleteStudentHandler)

	mux.HandleFunc("POST /execs/login", handlers.PatchStudentsHandler)	
	mux.HandleFunc("POST /execs/logout", handlers.PatchStudentsHandler)	
	mux.HandleFunc("POST /execs/forgotpassword", handlers.PatchStudentsHandler)	
	mux.HandleFunc("POST /execs/resetpassword/reset/{resetcode}", handlers.PatchStudentsHandler)	


	return mux

}
