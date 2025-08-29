package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, Root Route"))
	fmt.Println("Hello, Root Route")
}

func teachersHandler(w http.ResponseWriter, r *http.Request) {
	// GET /teachers/{id}

	switch r.Method {
	case http.MethodGet:
		fmt.Println(r.URL.Path)
		path := strings.TrimPrefix(r.URL.Path, "/teachers/")  // Extract the part after /teachers/
		fmt.Println("Path after /teachers/:", path)
		userID := strings.TrimSuffix(path, "/")  // browser might add a trailing slash
		fmt.Println("User ID:", userID)

		w.Write([]byte("GET Method on Teachers Route"))
	case http.MethodPost:
		w.Write([]byte("POST Method on Teachers Route"))
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

	http.HandleFunc("/", rootHandler)

	http.HandleFunc("/teachers/", teachersHandler)

	http.HandleFunc("/students/", studentsHandler)

	http.HandleFunc("/execs/", execsHandler)

	fmt.Println("Server is running on port", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Error starting the server: ", err)
	}

}
