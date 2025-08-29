package main

import (
	"encoding/json"
	"fmt"
	"io"
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
		// w.Write([]byte("Hello, Teachers Route"))
		// fmt.Println("Hello, Teachers Route")

		switch r.Method {
		case http.MethodGet:

			// Access the request details
			fmt.Println("Request Details:")
			fmt.Println("Form:", r.Form)
			fmt.Println("Header:", r.Header)
			fmt.Println("Context:", r.Context())
			fmt.Println("Content Length:", r.ContentLength)
			fmt.Println("Host:", r.Host)
			fmt.Println("Method:", r.Method)
			fmt.Println("URL:", r.URL)
			fmt.Println("Proto:", r.Proto)
			fmt.Println("Remote Address:", r.RemoteAddr)
			fmt.Println("Request URI:", r.RequestURI)
			fmt.Println("TLS:", r.TLS)
			fmt.Println("Trailer:", r.Trailer)
			fmt.Println("Transfer Encoding:", r.TransferEncoding)
			fmt.Println("User Agent:", r.UserAgent())
			fmt.Println("URL Port:", r.URL.Port())
			fmt.Println("URL Scheme:", r.URL.Scheme)

			w.Write([]byte("GET Method on Teachers Route"))
		case http.MethodPost:

			// Parse form data (necessary for x-www-form-urlencoded)
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Error parsing form data", http.StatusBadRequest)
				return
			}
			fmt.Println("Form Data:", r.Form)

			// Prepare response data
			response := make(map[string]interface{})
			for k, v := range r.Form {
				response[k] = v
			}
			fmt.Println("Processed Response Data:", response)

			// RAW Body (JSON)
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading body", http.StatusBadRequest)
				return
			}
			defer r.Body.Close()
			fmt.Println("Raw Body:", string(body))

			// if you expect JSON and want to unmarshal it into a struct
			var userInstance User
			err = json.Unmarshal(body, &userInstance)
			if err != nil {
				http.Error(w, "Error parsing JSON", http.StatusBadRequest)
				return
			}
			fmt.Println("Unmarshaled User Struct:", userInstance)
			fmt.Println("Received Name:", userInstance.Name)
			fmt.Println("Received Age:", userInstance.Age)
			fmt.Println("Received City:", userInstance.City)

			// Access the request details
			fmt.Println("Request Details:")
			fmt.Println("Body:", string(body)) // Already read earlier
			fmt.Println("Form:", r.Form)
			fmt.Println("Header:", r.Header)
			fmt.Println("Context:", r.Context())
			fmt.Println("Content Length:", r.ContentLength)
			fmt.Println("Host:", r.Host)
			fmt.Println("Method:", r.Method)
			fmt.Println("URL:", r.URL)
			fmt.Println("Proto:", r.Proto)
			fmt.Println("Remote Address:", r.RemoteAddr)
			fmt.Println("Request URI:", r.RequestURI)
			fmt.Println("TLS:", r.TLS)
			fmt.Println("Trailer:", r.Trailer)
			fmt.Println("Transfer Encoding:", r.TransferEncoding)
			fmt.Println("User Agent:", r.UserAgent())
			fmt.Println("URL Port:", r.URL.Port())
			fmt.Println("URL Scheme:", r.URL.Scheme)

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
		w.Write([]byte("Hello, Students Route"))
		fmt.Println("Hello, Students Route")
	})

	http.HandleFunc("/execs", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, Execs Route"))
		fmt.Println("Hello, Execs Route")
	})

	fmt.Println("Server is running on port", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Error starting the server: ", err)
	}

}
