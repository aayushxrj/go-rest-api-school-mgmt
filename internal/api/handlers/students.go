package handlers

import "net/http"

func StudentsHandler(w http.ResponseWriter, r *http.Request) {
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
