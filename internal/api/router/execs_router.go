package router

import (
	"net/http"

	"github.com/aayushxrj/go-rest-api-school-mgmt/internal/api/handlers"
)

func execsRouter(mux *http.ServeMux) {
	mux.HandleFunc("/execs/", handlers.ExecsHandler)
}
