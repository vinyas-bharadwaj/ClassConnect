package main

import (
	mw "ClassConnect/internal/api/middlewares"
	"ClassConnect/internal/api/routers"
	"ClassConnect/internal/repository/sqlconnect"
	"ClassConnect/pkg/utils"
	"os"

	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func init() {
	// Load .env file if it exists (for local development)
	_ = godotenv.Load()

	err := sqlconnect.InitDB()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	_, err = sqlconnect.ConnectDB()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func main() {
	port := ":" + os.Getenv("API_PORT")

	router := routers.Router()

	// Set a rate limiter of 5 requests per minute
	// rl := mw.NewRateLimiter(50, time.Minute)

	// Chaining all of our middlewares
	// Note that the first argument will be the innermost middleware and the last will be the outermost
	jwtMiddleware := mw.MiddlewareExcludePaths(mw.JWTMiddleware, "/execs/login/", "/execs/forgotPassword/", "/execs/resetPassword/", "/execs/")
	secureMux := utils.ApplyMiddlewares(router, mw.SecurityHeaders, mw.Compress, jwtMiddleware, mw.ResponseTime, mw.Cors)

	// Create custom server
	server := &http.Server{
		Addr:    port,
		Handler: secureMux,
	}

	fmt.Println("Server running on port:", port, "(HTTPS)")
	err := server.ListenAndServeTLS("/app/server.crt", "/app/server.key")
	if err != nil {
		log.Fatalln("Error starting new server: ", err)
	}
}
