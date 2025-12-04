package routers

import (
	"ClassConnect/internal/api/handlers"
	"ClassConnect/internal/repository/sqlconnect"
	"log"
	"net/http"
)

func teachersRouter() *http.ServeMux {
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		log.Fatal("Error:", err)
		return nil
	}

	mux := http.NewServeMux()
	teacherHandler := handlers.NewTeacherHandler(db)

	// Teacher routes
	mux.HandleFunc("GET /teachers/", teacherHandler.GetTeachersHandler)
	mux.HandleFunc("POST /teachers/", teacherHandler.CreateTeachersHandler)

	mux.HandleFunc("GET /teachers/{id}", teacherHandler.GetTeacherByIdhandler)
	mux.HandleFunc("PUT /teachers/{id}", teacherHandler.UpdateTeachersHandler)
	mux.HandleFunc("DELETE /teachers/{id}", teacherHandler.DeleteTeachersHandler)

	mux.HandleFunc("GET /teachers/{id}/students", teacherHandler.GetStudentsByTeacherId)

	return mux
}
