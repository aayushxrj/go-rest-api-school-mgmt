package handlers

import "net/http"

func ExecsHandler(w http.ResponseWriter, r *http.Request) {
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
