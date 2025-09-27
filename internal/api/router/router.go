package router

import (
	"net/http"
)

func MainRouter() *http.ServeMux {

	mux := http.NewServeMux()

	TeachersRouter(mux)
	StudentsRouter(mux)
	ExecsRouter(mux)

	return mux

}
