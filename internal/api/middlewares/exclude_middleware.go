package middlewares

import (
	"net/http"
	"strings"
)

func MiddlewareExcludePaths(middleware func(http.Handler) http.Handler, excludedPaths ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			for _, path := range excludedPaths {
				if strings.HasPrefix(r.URL.Path, path) {
					// Skip the JWT middleware if the specified path is present
					next.ServeHTTP(w, r)
					return
				}
			}

			middleware(next).ServeHTTP(w, r)
		})
	}
}
