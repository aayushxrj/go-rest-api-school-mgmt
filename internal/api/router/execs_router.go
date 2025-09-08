package router

import (
	"net/http"

	"github.com/aayushxrj/go-rest-api-school-mgmt/internal/api/handlers"
)

func execsRouter(mux *http.ServeMux) {
	
	mux.HandleFunc("GET /execs", handlers.GetExecsHandler)
	mux.HandleFunc("POST /execs", handlers.AddExecHandler)
	mux.HandleFunc("PATCH /execs", handlers.PatchExecsHandler)
	
	mux.HandleFunc("GET /execs/{id}", handlers.GetOneExecHandler)
	mux.HandleFunc("PATCH /execs/{id}", handlers.PatchOneExecHandler)
	mux.HandleFunc("DELETE /execs/{id}", handlers.DeleteOneExecHandler)
	// mux.HandleFunc("POST /execs/{id}/updatepassword", handlers.ExecsHandler)
	
	mux.HandleFunc("POST /execs/login", handlers.LoginHandler)
	// mux.HandleFunc("POST /execs/logout", handlers.ExecsHandler)
	// mux.HandleFunc("POST /execs/forotpassword", handlers.ExecsHandler)
	// mux.HandleFunc("POST /execs/resetpassword/reset/{resetcode}", handlers.ExecsHandler)
}
