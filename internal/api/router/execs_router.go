package router

import (
	"net/http"

	"github.com/aayushxrj/go-rest-api-school-mgmt/internal/api/handlers"
)

func execsRouter(mux *http.ServeMux) {
	
	mux.HandleFunc("GET /execs", handlers.ExecsHandler)
	mux.HandleFunc("POST /execs", handlers.ExecsHandler)
	mux.HandleFunc("PATCH /execs", handlers.ExecsHandler)
	
	mux.HandleFunc("GET /execs/{id}", handlers.ExecsHandler)
	mux.HandleFunc("PATCH /execs/{id}", handlers.ExecsHandler)
	mux.HandleFunc("DELETE /execs/{id}", handlers.ExecsHandler)
	mux.HandleFunc("POST /execs/{id}/updatepassword", handlers.ExecsHandler)
	
	mux.HandleFunc("POST /execs/login", handlers.ExecsHandler)
	mux.HandleFunc("POST /execs/logout", handlers.ExecsHandler)
	mux.HandleFunc("POST /execs/forotpassword", handlers.ExecsHandler)
	mux.HandleFunc("POST /execs/resetpassword/reset/{resetcode}", handlers.ExecsHandler)
}
