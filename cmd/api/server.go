// @title School Management API
// @version 1.0
// @description REST API for managing school data.
// @schemes https
// @host localhost:3000
// @BasePath /
package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"

	// "time"

	mw "github.com/aayushxrj/go-rest-api-school-mgmt/internal/api/middlewares"
	"github.com/aayushxrj/go-rest-api-school-mgmt/internal/api/router"
	"github.com/aayushxrj/go-rest-api-school-mgmt/internal/repository/sqlconnect"
	"github.com/aayushxrj/go-rest-api-school-mgmt/pkg/utils"
	"github.com/joho/godotenv"
	"golang.org/x/net/http2"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	_, err = sqlconnect.ConnectDB()
	if err != nil {
		// log.Fatal("Error connecting to database:", err)
		utils.ErrorHandler(err, "Error connecting to database")
		return
	}

	port := fmt.Sprintf(":%s", os.Getenv("API_PORT"))

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

	jwtMiddleware := mw.MiddlewaresExcludePaths(mw.JWTMiddleware,
		"/swagger",
		"/execs/login",
		"/execs/forgotpassword",
		"/execs/resetpassword/reset")

	// proper ordering of middlewares
	// example: Cors -> Rate Limiter -> Response Time -> Security Headers -> Compression -> HPP -> Actual Handler
	// secureMux := mw.Cors(rl.Middleware(mw.ResponseTimeMiddleware(mw.SecurityHeaders(mw.Compression(mw.Hpp(hppOptions)(mux))))))
	secureMux := utils.ApplyMiddlewares(router,
		// mw.Hpp(hppOptions),
		// mw.Compression,
		mw.SecurityHeaders,
		jwtMiddleware,
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

	err = server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}

}
