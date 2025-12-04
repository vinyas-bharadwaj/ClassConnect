package handlers

import (
	"ClassConnect/internal/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type StudentHandler struct {
	db *sql.DB
}

func NewStudentHandler(db *sql.DB) *StudentHandler {
	return &StudentHandler{db: db}
}

func (h *StudentHandler) GetStudentsHandler(w http.ResponseWriter, _ *http.Request) {
	rows, err := h.db.Query("SELECT * FROM students;")
	if err != nil {
		http.Error(w, "Error retrieving all the students", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	studentList := make([]models.Student, 0)
	for rows.Next() {
		var student models.Student
		err = rows.Scan(&student.Id, &student.FirstName, &student.LastName, &student.Email, &student.Class)
		if err != nil {
			http.Error(w, "Error scanning database results", http.StatusInternalServerError)
			return
		}
		studentList = append(studentList, student)
	}
	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Student `json:"data"`
	}{
		Status: "success",
		Count:  len(studentList),
		Data:   studentList,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *StudentHandler) GetStudentByIdHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(err)
		return
	}

	var student models.Student
	err = h.db.QueryRow("SELECT * FROM students WHERE id = ?", id).Scan(
		&student.Id,
		&student.FirstName,
		&student.LastName,
		&student.Email,
		&student.Class,
	)
	if err != nil {
		http.Error(w, "Student with that ID does not exist in the database!", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

func (h *StudentHandler) CreateStudentsHandler(w http.ResponseWriter, r *http.Request) {
	var newStudents []models.Student
	err := json.NewDecoder(r.Body).Decode(&newStudents)
	if err != nil {
		http.Error(w, "Invalid Request body", http.StatusBadRequest)
		return
	}

	stmt, err := h.db.Prepare("INSERT INTO students(first_name, last_name, email, class) VALUES(?,?,?,?)")
	if err != nil {
		http.Error(w, "Error in preparing SQL query", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	addedStudents := make([]models.Student, len(newStudents))
	for i, newStudent := range newStudents {
		res, err := stmt.Exec(newStudent.FirstName, newStudent.LastName, newStudent.Email, newStudent.Class)
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
		newStudent.Id = int(lastId)
		addedStudents[i] = newStudent
	}

	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Student `json:"data"`
	}{
		Status: "success",
		Count:  len(addedStudents),
		Data:   addedStudents,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func (h *StudentHandler) DeleteStudentsHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := h.db.Exec("DELETE FROM students WHERE id = ?", id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error deleting the student", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Println(err)
		http.Error(w, "Error deleting the student", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		log.Println(err)
		http.Error(w, "The student does not exist", http.StatusInternalServerError)
		return
	}

	response := struct {
		Status string `json:"status"`
		Id     int    `json:"id"`
	}{
		Status: "Successfully deleted the student",
		Id:     id,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *StudentHandler) UpdateStudentsHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Student ID", http.StatusBadGateway)
		return
	}

	var updatedStudent models.Student
	err = json.NewDecoder(r.Body).Decode(&updatedStudent)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid request payload", http.StatusBadGateway)
		return
	}

	var existingStudent models.Student
	err = h.db.QueryRow("SELECT * FROM students WHERE id = ?", id).Scan(
		&existingStudent.Id,
		&existingStudent.FirstName,
		&existingStudent.LastName,
		&existingStudent.Email,
		&existingStudent.Class,
	)
	// In case there are no rows, an error is still thrown
	if err == sql.ErrNoRows {
		http.Error(w, "Student with the given ID not found!", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Unable to retrieve data", http.StatusInternalServerError)
		return
	}

	updatedStudent.Id = existingStudent.Id
	_, err = h.db.Exec("UPDATE students SET first_name = ?, last_name = ?, email = ?, class = ? WHERE id = ?",
		updatedStudent.FirstName,
		updatedStudent.LastName,
		updatedStudent.Email,
		updatedStudent.Class,
		id,
	)

	if err != nil {
		http.Error(w, "Error updating the students details", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedStudent)
}
