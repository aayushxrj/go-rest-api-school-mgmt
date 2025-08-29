package main

import (
	"fmt"
	"log"
	"net/http"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	City string `json:"city"`
}

func main() {

	port := ":3000"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, Root Route"))
		fmt.Println("Hello, Root Route")
	})

	http.HandleFunc("/teachers", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
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
	})

	http.HandleFunc("/students", func(w http.ResponseWriter, r *http.Request) {
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
	})

	http.HandleFunc("/execs", func(w http.ResponseWriter, r *http.Request) {
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
	})

	fmt.Println("Server is running on port", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Error starting the server: ", err)
	}

}
