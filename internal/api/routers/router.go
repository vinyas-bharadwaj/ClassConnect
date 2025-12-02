package routers

import (
	"ClassConnect/internal/api/handlers"
	"ClassConnect/internal/repository/sqlconnect"
	"log"

	"net/http"
)

func Router() *http.ServeMux {
	mux := http.NewServeMux()
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		log.Fatal("Error:", err)
		return nil
	}

	teacherHandler := handlers.NewTeacherHandler(db)

	mux.HandleFunc("GET /teachers/", teacherHandler.GetTeachersHandler)
	mux.HandleFunc("POST /teachers/", teacherHandler.CreateTeachersHandler)
	mux.HandleFunc("PUT /teachers/{id}", teacherHandler.UpdateTeachersHandler)
	mux.HandleFunc("DELETE /teachers/{id}", teacherHandler.DeleteTeachersHandler)

	return mux
}
