package routers

import (
	"ClassConnect/internal/api/handlers"
	"ClassConnect/internal/repository/sqlconnect"
	"log"
	"net/http"
)

func execsRouter() *http.ServeMux {
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		log.Fatal("Error:", err)
		return nil
	}

	mux := http.NewServeMux()
	execsHandler := handlers.NewExecsHandler(db)

	// Execs routes
	mux.HandleFunc("GET /execs/", execsHandler.GetExecsHandler)
	mux.HandleFunc("POST /execs/", execsHandler.CreateExecsHandler)

	mux.HandleFunc("GET /execs/{id}", execsHandler.GetExecByIdHandler)
	mux.HandleFunc("PUT /execs/{id}", execsHandler.UpdateExecsHandler)
	mux.HandleFunc("DELETE /execs/{id}", execsHandler.DeleteExecsHandler)

	mux.HandleFunc("POST /execs/login/", execsHandler.LoginHandler)
	mux.HandleFunc("POST /execs/logout/", execsHandler.LogoutHandler)
	mux.HandleFunc("PATCH /execs/updatePassword/{id}", execsHandler.UpdatePasswordHandler)
	mux.HandleFunc("POST /execs/forgotPassword/", execsHandler.ForgotPasswordHandler)
	mux.HandleFunc("POST /execs/resetPassword/{resetCode}", execsHandler.ResetPasswordHandler)

	return mux
}
