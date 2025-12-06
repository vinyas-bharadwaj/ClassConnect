package handlers

import (
	"ClassConnect/internal/models"
	"ClassConnect/pkg/utils"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-mail/mail/v2"
	"golang.org/x/crypto/argon2"
)

type ExecsHandler struct {
	db *sql.DB
}

func NewExecsHandler(db *sql.DB) *ExecsHandler {
	return &ExecsHandler{db: db}
}

func (h *ExecsHandler) GetExecsHandler(w http.ResponseWriter, _ *http.Request) {
	rows, err := h.db.Query("SELECT id, first_name, last_name, email, username, password, password_changed_at, user_created_at, password_reset_token, inactive_status, role FROM execs;")
	if err != nil {
		log.Println("Query error:", err)
		http.Error(w, "Error retrieving all the execs", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	execsList := make([]models.Exec, 0)
	for rows.Next() {
		var exec models.Exec
		err = rows.Scan(&exec.Id, &exec.FirstName, &exec.LastName, &exec.Email, &exec.Username, &exec.Password, &exec.PasswordChangedAt, &exec.UserCreatedAt, &exec.PasswordResetCode, &exec.InactiveStatus, &exec.Role)
		if err != nil {
			log.Println("Scan error:", err)
			http.Error(w, "Error scanning database results", http.StatusInternalServerError)
			return
		}
		execsList = append(execsList, exec)
	}
	response := struct {
		Status string        `json:"status"`
		Count  int           `json:"count"`
		Data   []models.Exec `json:"data"`
	}{
		Status: "success",
		Count:  len(execsList),
		Data:   execsList,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *ExecsHandler) GetExecByIdHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(err)
		return
	}

	var exec models.Exec
	err = h.db.QueryRow("SELECT id, first_name, last_name, email, username, password, password_changed_at, user_created_at, password_reset_token, inactive_status, role FROM execs WHERE id = ?", id).Scan(
		&exec.Id,
		&exec.FirstName,
		&exec.LastName,
		&exec.Email,
		&exec.Username,
		&exec.Password,
		&exec.PasswordChangedAt,
		&exec.UserCreatedAt,
		&exec.PasswordResetCode,
		&exec.InactiveStatus,
		&exec.Role,
	)
	if err != nil {
		log.Println("Query error:", err)
		http.Error(w, "Exec with that ID does not exist in the database!", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exec)
}

func (h *ExecsHandler) CreateExecsHandler(w http.ResponseWriter, r *http.Request) {
	var newExecs []models.Exec
	err := json.NewDecoder(r.Body).Decode(&newExecs)
	if err != nil {
		http.Error(w, "Invalid Request body", http.StatusBadRequest)
		return
	}

	stmt, err := h.db.Prepare("INSERT INTO execs(first_name, last_name, email, username, password, password_changed_at, password_reset_token, inactive_status, role) VALUES(?,?,?,?,?,?,?,?,?)")
	if err != nil {
		log.Println("Prepare error:", err)
		http.Error(w, "Error in preparing SQL query", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	addedExecs := make([]models.Exec, len(newExecs))
	for i, newExec := range newExecs {
		if newExec.Password == "" {
			http.Error(w, "Password field cannot be empty", http.StatusBadRequest)
			return
		}

		salt := make([]byte, 16)
		_, err := rand.Read(salt)
		if err != nil {
			http.Error(w, "Error generating the password hash", http.StatusInternalServerError)
			return
		}

		hash := argon2.IDKey([]byte(newExec.Password), salt, 1, 64*1024, 4, 32)
		saltBase64 := base64.StdEncoding.EncodeToString(salt)
		hashBase64 := base64.StdEncoding.EncodeToString(hash)

		encodedHash := fmt.Sprintf("%s.%s", saltBase64, hashBase64)
		newExec.Password = encodedHash

		res, err := stmt.Exec(
			newExec.FirstName,
			newExec.LastName,
			newExec.Email,
			newExec.Username,
			newExec.Password,
			newExec.PasswordChangedAt,
			newExec.PasswordResetCode,
			newExec.InactiveStatus,
			newExec.Role,
		)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error inserting data into the database", http.StatusInternalServerError)
			return
		}
		lastId, err := res.LastInsertId()
		if err != nil {
			log.Println(err)
			http.Error(w, "Error getting last insert id", http.StatusInternalServerError)
			return
		}
		newExec.Id = int(lastId)
		addedExecs[i] = newExec
	}

	response := struct {
		Status string        `json:"status"`
		Count  int           `json:"count"`
		Data   []models.Exec `json:"data"`
	}{
		Status: "success",
		Count:  len(addedExecs),
		Data:   addedExecs,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func (h *ExecsHandler) DeleteExecsHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := h.db.Exec("DELETE FROM execs WHERE id = ?", id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error deleting the exec", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Println(err)
		http.Error(w, "Error deleting the exec", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		log.Println(err)
		http.Error(w, "The exec does not exist", http.StatusInternalServerError)
		return
	}

	response := struct {
		Status string `json:"status"`
		Id     int    `json:"id"`
	}{
		Status: "Successfully deleted the exec",
		Id:     id,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *ExecsHandler) UpdateExecsHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Exec ID", http.StatusBadGateway)
		return
	}

	var updatedExec models.Exec
	err = json.NewDecoder(r.Body).Decode(&updatedExec)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid request payload", http.StatusBadGateway)
		return
	}

	var existingExec models.Exec
	err = h.db.QueryRow("SELECT id, first_name, last_name, email, username, password, password_changed_at, user_created_at, password_reset_token, inactive_status, role FROM execs WHERE id = ?", id).Scan(
		&existingExec.Id,
		&existingExec.FirstName,
		&existingExec.LastName,
		&existingExec.Email,
		&existingExec.Username,
		&existingExec.Password,
		&existingExec.PasswordChangedAt,
		&existingExec.UserCreatedAt,
		&existingExec.PasswordResetCode,
		&existingExec.InactiveStatus,
		&existingExec.Role,
	)
	// In case there are no rows, an error is still thrown
	if err == sql.ErrNoRows {
		http.Error(w, "Exec with the given ID not found!", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Unable to retrieve data", http.StatusInternalServerError)
		return
	}

	updatedExec.Id = existingExec.Id
	_, err = h.db.Exec("UPDATE execs SET first_name = ?, last_name = ?, email = ?, username = ?, password = ?, password_changed_at = ?, password_reset_token = ?, inactive_status = ?, role = ? WHERE id = ?",
		updatedExec.FirstName,
		updatedExec.LastName,
		updatedExec.Email,
		updatedExec.Username,
		updatedExec.Password,
		updatedExec.PasswordChangedAt,
		updatedExec.PasswordResetCode,
		updatedExec.InactiveStatus,
		updatedExec.Role,
		id,
	)

	if err != nil {
		http.Error(w, "Error updating the execs details", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedExec)
}

func (h *ExecsHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req models.Exec
	// Data validation
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.Username == "" || req.Password == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Search for user if they exist in the database
	var user models.Exec
	err = h.db.QueryRow(
		"SELECT id, first_name, last_name, email, username, password, inactive_status, role FROM execs WHERE username = ?",
		req.Username,
	).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Username, &user.Password, &user.InactiveStatus, &user.Role)

	if err != nil {
		http.Error(w, "Error locating the user in the database", http.StatusNotFound)
		return
	}

	// Check if user is active
	if user.InactiveStatus {
		http.Error(w, "Account is inactive", http.StatusForbidden)
		return
	}

	// Verify password
	err = utils.VerifyPassword(req.Password, user.Password)
	if err != nil {
		http.Error(w, "Incorrect password", http.StatusForbidden)
		return
	}

	// Generate JWT token
	tokenString, err := utils.SignToken(strconv.Itoa(user.Id), user.Username, user.Role)
	if err != nil {
		http.Error(w, "Error generating JWT token", http.StatusInternalServerError)
		return
	}

	// Send token as a response or as a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "Bearer",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(24 * time.Hour),
		SameSite: http.SameSiteStrictMode,
	})

	// Response body
	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Token string `json:"token"`
	}{
		Token: tokenString,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *ExecsHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clear the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "Bearer",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Unix(0, 0),
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Message string `json:"message"`
	}{
		Message: "Logged out successfully",
	}

	json.NewEncoder(w).Encode(response)
}

func (h *ExecsHandler) UpdatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	userId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid exec id", http.StatusBadRequest)
		return
	}

	var req models.UpdatePasswordRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	r.Body.Close()

	if req.CurrentPassword == "" || req.NewPassword == "" {
		http.Error(w, "Password cannot be blank", http.StatusBadRequest)
		return
	}

	var username string
	var userPassword string
	var role string
	err = h.db.QueryRow("SELECT username, password, role FROM execs WHERE id = ?", userId).Scan(&username, &userPassword, &role)
	if err != nil {
		http.Error(w, "User with the ID does not exist", http.StatusNotFound)
		return
	}

	err = utils.VerifyPassword(req.CurrentPassword, userPassword)
	if err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// Hash the new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		log.Println("Hash error:", err)
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Update the password in the database
	_, err = h.db.Exec("UPDATE execs SET password = ?, password_changed_at = CURRENT_TIMESTAMP WHERE id = ?", hashedPassword, userId)
	if err != nil {
		log.Println("Update error:", err)
		http.Error(w, "Error updating password", http.StatusInternalServerError)
		return
	}

	tokenString, err := utils.SignToken(strconv.Itoa(userId), username, role)
	if err != nil {
		http.Error(w, "Error generating JWT token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "Bearer",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(24 * time.Hour),
		SameSite: http.SameSiteStrictMode,
	})

	response := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Token   string `json:"token"`
	}{
		Status:  "success",
		Message: "Password updated successfully",
		Token:   tokenString,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *ExecsHandler) ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("ForgotPasswordHandler called")

	var req struct {
		Email string `json:"email"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("Decode error:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	r.Body.Close()

	if req.Email == "" {
		log.Println("Email field is blank")
		http.Error(w, "Email field should not be blank", http.StatusBadRequest)
		return
	}

	log.Printf("Looking up user with email: %s\n", req.Email)

	var exec models.Exec
	err = h.db.QueryRow("SELECT id FROM execs WHERE email = ?", req.Email).Scan(&exec.Id)
	if err != nil {
		log.Println("User lookup error:", err)
		http.Error(w, "User with that email does not exist", http.StatusNotFound)
		return
	}

	log.Printf("User found with ID: %d\n", exec.Id)

	durationStr := os.Getenv("RESET_TOKEN_EXP_DURATION")
	if durationStr == "" {
		log.Println("RESET_TOKEN_EXP_DURATION not set, using default 15 minutes")
		durationStr = "15"
	}

	duration, err := strconv.Atoi(durationStr)
	if err != nil {
		log.Println("Duration parse error:", err)
		http.Error(w, "Error loading in the environment variable", http.StatusInternalServerError)
		return
	}

	tokenBytes := make([]byte, 32)
	_, err = rand.Read(tokenBytes)
	if err != nil {
		log.Println("Token generation error:", err)
		http.Error(w, "Failed to send password reset email", http.StatusInternalServerError)
		return
	}

	token := hex.EncodeToString(tokenBytes)
	hashedToken := sha256.Sum256(tokenBytes)
	log.Println("Generated token (save this for testing):", token)

	hashedTokenString := hex.EncodeToString(hashedToken[:])

	// Note: Using password_reset_token column from schema
	_, err = h.db.Exec("UPDATE execs SET password_reset_token = ? WHERE id = ?",
		hashedTokenString,
		exec.Id,
	)
	if err != nil {
		log.Println("Database update error:", err)
		http.Error(w, "Failed to save password reset token", http.StatusInternalServerError)
		return
	}

	log.Println("Password reset token saved to database")

	// Send reset email
	resetURL := fmt.Sprintf("http://localhost:3000/execs/resetPassword/%s", token)
	message := fmt.Sprintf("Forgot your password? Reset your password using the following link: \n%s\n\nIf you didn't request a password reset, please ignore this email. This link is only valid for %d minutes", resetURL, duration)

	log.Printf("Sending email to: %s\n", req.Email)

	m := mail.NewMessage()
	m.SetHeader("From", "schooladmin@school.com")
	m.SetHeader("To", req.Email)
	m.SetHeader("Subject", "Your password reset link")
	m.SetBody("text/plain", message)

	d := mail.NewDialer("localhost", 1025, "", "")
	err = d.DialAndSend(m)
	if err != nil {
		log.Println("Email send error:", err)
		http.Error(w, "Failed to send the email", http.StatusInternalServerError)
		return
	}

	log.Println("Email sent successfully")

	// Respond with success message
	response := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}{
		Status:  "success",
		Message: fmt.Sprintf("Password reset link sent to %s", req.Email),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
