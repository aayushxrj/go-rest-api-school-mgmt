package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	mw "github.com/aayushxrj/go-rest-api-school-mgmt/internal/api/middlewares"
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

	mux := http.NewServeMux()

	mux.HandleFunc("/", rootHandler)

	mux.HandleFunc("/teachers/", teachersHandler)

	mux.HandleFunc("/students/", studentsHandler)

	mux.HandleFunc("/execs/", execsHandler)

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// rate limiter
	rl := mw.NewRateLimiter(5, 1*time.Minute)

	hppOptions := mw.HPPOptions{
		CheckQuery:                  true,
		CheckBody:                   true,
		CheckBodyOnlyForContentType: "application/x-www-form-urlencoded",
		Whitelist:                   []string{"sortBy", "sortOrder", "name", "age", "class"},
	}

	// proper ordering of middlewares
	// example: Cors -> Rate Limiter -> Response Time -> Security Headers -> Compression -> HPP -> Actual Handler
	// secureMux := mw.Cors(rl.Middleware(mw.ResponseTimeMiddleware(mw.SecurityHeaders(mw.Compression(mw.Hpp(hppOptions)(mux))))))
	secureMux := applyMiddlewares(mux,
		mw.Hpp(hppOptions),
		mw.Compression,
		mw.SecurityHeaders,
		mw.ResponseTimeMiddleware,
		rl.Middleware,
		mw.Cors,
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
