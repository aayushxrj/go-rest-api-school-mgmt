package router

import (
	"net/http"

	"github.com/aayushxrj/go-rest-api-school-mgmt/internal/api/handlers"
	_ "github.com/aayushxrj/go-rest-api-school-mgmt/docs" 
    httpSwagger "github.com/swaggo/http-swagger"
)

func Router() *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.RootHandler)

	// Swagger UI route
    mux.Handle("/swagger/", httpSwagger.WrapHandler)

	studentsRouter(mux)
	teachersRouter(mux)
	execsRouter(mux)

	return mux
}
