package middlewares

import (
	"net/http"
)

var allowedOrigins = []string{
	"http://localhost:3000",
	"http://localhost:5173",
}

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if !isAllowedOrigin(origin) {
			http.Error(w, "Not allowed by CORS", http.StatusForbidden)
			return
		} else {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		// Set the cors headers
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// The options http method is just a pre-flight check performed by the browsers
		if r.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isAllowedOrigin(origin string) bool {
	for _, allowedOrigin := range allowedOrigins {
		if origin == allowedOrigin {
			return true
		}
	}
	return false
}
