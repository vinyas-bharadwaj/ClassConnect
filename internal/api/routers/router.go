package routers

import (
	"ClassConnect/internal/api/handlers"
	"net/http"
)

func Router() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/teachers/", handlers.TeachersHandler)

	return mux
}
