package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	// "time"

	mw "github.com/aayushxrj/go-rest-api-school-mgmt/internal/api/middlewares"
	"github.com/aayushxrj/go-rest-api-school-mgmt/internal/api/router"
	"github.com/aayushxrj/go-rest-api-school-mgmt/pkg/utils"
	"golang.org/x/net/http2"
)

func main() {

	port := ":3000"

	cert := "cmd/api/cert.pem"
	key := "cmd/api/key.pem"

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	router := router.Router()

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
	secureMux := utils.ApplyMiddlewares(router,
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
