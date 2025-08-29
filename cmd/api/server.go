package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/http2"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, Root Route"))
	fmt.Println("Hello, Root Route")
}

func teachersHandler(w http.ResponseWriter, r *http.Request) {
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

	http.HandleFunc("/", rootHandler)

	http.HandleFunc("/teachers/", teachersHandler)

	http.HandleFunc("/students/", studentsHandler)

	http.HandleFunc("/execs/", execsHandler)

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// Create custom server
	server := &http.Server{
		Addr:      port,
		Handler:   nil, // use default mux
		TLSConfig: tlsConfig,
	}

	http2.ConfigureServer(server, &http2.Server{})

	fmt.Println("Server Listening on port", port)

	err := server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}

	// fmt.Println("Server is running on port", port)
	// err := http.ListenAndServe(port, nil)
	// if err != nil {
	// 	log.Fatal("Error starting the server: ", err)
	// }

}
