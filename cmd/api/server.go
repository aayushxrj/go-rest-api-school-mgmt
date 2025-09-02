package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	// "time"

	mw "github.com/aayushxrj/go-rest-api-school-mgmt/internal/api/middlewares"
	"golang.org/x/net/http2"
)

type Teacher struct {
	ID        int    `json:"id,omitempty`
	FirstName string `json :"first_name,omitempty"`
	LastName  string `json :"last_name,omitempty"`
	Class     string    `json :"class,omitempty"`
	Subject   string `json :"subject,omitempty"`
}

var (
	teachers = make(map[int]Teacher)
	mutex    = &sync.Mutex{} // use in the post method
	nextID   = 1
)

func init() {
	teachers[nextID] = Teacher{ID: nextID, FirstName: "John", LastName: "Doe", Class: "10", Subject: "Math"}
	nextID++
	teachers[nextID] = Teacher{ID: nextID, FirstName: "Jane", LastName: "Smith", Class: "9", Subject: "Science"}
	nextID++
	teachers[nextID] = Teacher{ID: nextID, FirstName: "Jane", LastName: "Doe", Class: "11", Subject: "Biology"}
	nextID++
}

func getTeachersHandler(w http.ResponseWriter, r *http.Request) {

	// Path Parameters
	pathStr := strings.TrimPrefix(r.URL.Path, "/teachers/")
	idStr := strings.TrimSuffix(pathStr, "/")
	// fmt.Println("ID String:", idStr)

	// Query Parameters
	if idStr == "" {
		first_name := r.URL.Query().Get("first_name")
		last_name := r.URL.Query().Get("last_name")

		teacherList := make([]Teacher, 0, len(teachers))
		for _, teacher := range teachers {
			if (first_name == "" || first_name == teacher.FirstName) && (last_name == "" || last_name == teacher.LastName) {
				teacherList = append(teacherList, teacher)
			}
		}

		response := struct {
			Status string    `json:"status"`
			Count  int       `json:"count"`
			Data   []Teacher `json:"data"`
		}{
			Status: "success",
			Count:  len(teacherList),
			Data:   teacherList,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Handle Path Parameter
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Teacher ID", http.StatusBadRequest)
		return
	}

	teacher, exists := teachers[id]
	if !exists {
		http.Error(w, "Teacher Not Found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teacher)
}

func addTeacherHadnler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	var newTeachers []Teacher
	err := json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	addedTeachers := make([]Teacher, len(newTeachers))

	for i, newTeacher := range newTeachers {
		newTeacher.ID = nextID
		teachers[nextID] = newTeacher
		addedTeachers[i] = newTeacher
		nextID++
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string    `json:"status"`
		Count  int       `json:"count"`
		Data   []Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(addedTeachers),
		Data:   addedTeachers,
	}

	json.NewEncoder(w).Encode(response)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, Root Route"))
	fmt.Println("Hello, Root Route")
}

func teachersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTeachersHandler(w, r)
	case http.MethodPost:
		addTeacherHadnler(w, r)
	case http.MethodPut:
		w.Write([]byte("PUT Method on Teachers Route"))
	case http.MethodPatch:
		w.Write([]byte("PATCH Method on Teachers Route"))
	case http.MethodDelete:
		w.Write([]byte("DELETE Method on Teachers Route"))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed"))
	}
}

func studentsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("GET Method on Students Route"))
	case http.MethodPost:
		w.Write([]byte("POST Method on Students Route"))
	case http.MethodPut:
		w.Write([]byte("PUT Method on Students Route"))
	case http.MethodPatch:
		w.Write([]byte("PATCH Method on Students Route"))
	case http.MethodDelete:
		w.Write([]byte("DELETE Method on Students Route"))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed"))
	}
}

func execsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("GET Method on Execs Route"))
	case http.MethodPost:
		w.Write([]byte("POST Method on Execs Route"))
	case http.MethodPut:
		w.Write([]byte("PUT Method on Execs Route"))
	case http.MethodPatch:
		w.Write([]byte("PATCH Method on Execs Route"))
	case http.MethodDelete:
		w.Write([]byte("DELETE Method on Execs Route"))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed"))
	}
}

func main() {

	port := ":3000"

	cert := "cmd/api/cert.pem"
	key := "cmd/api/key.pem"

	mux := http.NewServeMux()

	mux.HandleFunc("/", rootHandler)

	mux.HandleFunc("/teachers/", teachersHandler)

	mux.HandleFunc("/students/", studentsHandler)

	mux.HandleFunc("/execs/", execsHandler)

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// rate limiter
	// rl := mw.NewRateLimiter(5, 1*time.Minute)

	// hppOptions := mw.HPPOptions{
	// 	CheckQuery:                  true,
	// 	CheckBody:                   true,
	// 	CheckBodyOnlyForContentType: "application/x-www-form-urlencoded",
	// 	Whitelist:                   []string{"sortBy", "sortOrder", "name", "age", "class"},
	// }

	// proper ordering of middlewares
	// example: Cors -> Rate Limiter -> Response Time -> Security Headers -> Compression -> HPP -> Actual Handler
	// secureMux := mw.Cors(rl.Middleware(mw.ResponseTimeMiddleware(mw.SecurityHeaders(mw.Compression(mw.Hpp(hppOptions)(mux))))))
	secureMux := applyMiddlewares(mux,
		// mw.Hpp(hppOptions),
		// mw.Compression,
		mw.SecurityHeaders,
		// mw.ResponseTimeMiddleware,
		// rl.Middleware,
		// mw.Cors,
	)

	// Create custom server
	server := &http.Server{
		Addr:      port,
		Handler:   secureMux,
		TLSConfig: tlsConfig,
	}

	http2.ConfigureServer(server, &http2.Server{})

	fmt.Println("Server Listening on port", port)

	err := server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}

}

// Middleware is a function that wraps an http.Handler with additional functionality
type Middleware func(http.Handler) http.Handler

func applyMiddlewares(handler http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}
