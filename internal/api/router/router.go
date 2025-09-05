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

	mux.HandleFunc("GET /teachers/", handlers.GetTeachersHandler)
	mux.HandleFunc("POST /teachers/", handlers.AddTeacherHandler)
	mux.HandleFunc("PATCH /teachers/", handlers.PatchTeachersHandler)
	mux.HandleFunc("DELETE /teachers/", handlers.DeleteTeachersHandler)
	
	mux.HandleFunc("GET /teachers/{id}", handlers.GetOneTeacherHandler)
	mux.HandleFunc("PUT /teachers/{id}", handlers.UpdateTeacherHandler)
	mux.HandleFunc("PATCH /teachers/{id}", handlers.PatchOneTeacherHandler)
	mux.HandleFunc("DELETE /teachers/{id}", handlers.DeleteOneTeacherHandler)

	mux.HandleFunc("/students/", handlers.StudentsHandler)

	mux.HandleFunc("/execs/", handlers.ExecsHandler)

	return mux
}
