package routers

import (
	"ClassConnect/internal/api/handlers"
	"ClassConnect/internal/repository/sqlconnect"
	"log"
	"net/http"
)

func studentsRouter() *http.ServeMux {
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		log.Fatal("Error:", err)
		return nil
	}

	mux := http.NewServeMux()
	studentHandler := handlers.NewStudentHandler(db)

	// Student routes
	mux.HandleFunc("GET /students/", studentHandler.GetStudentsHandler)
	mux.HandleFunc("POST /students/", studentHandler.CreateStudentsHandler)

	mux.HandleFunc("GET /students/{id}", studentHandler.GetStudentByIdHandler)
	mux.HandleFunc("PUT /students/{id}", studentHandler.UpdateStudentsHandler)
	mux.HandleFunc("DELETE /students/{id}", studentHandler.DeleteStudentsHandler)

	return mux
}
