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
	"time"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	err = sqlconnect.InitDB()
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
	rl := mw.NewRateLimiter(50, time.Minute)

	// Chaining all of our middlewares
	// Note that the first argument will be the innermost middleware and the last will be the outermost
	secureMux := utils.ApplyMiddlewares(router, mw.Compress, mw.SecurityHeaders, mw.ResponseTime, rl.Middleware, mw.Cors)

	// Create custom server
	server := &http.Server{
		Addr:    port,
		Handler: secureMux,
	}

	fmt.Println("Server running on port:", port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln("Error starting new server: ", err)
	}
}
