package handlers

import (
	"ClassConnect/internal/models"
	"ClassConnect/pkg/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type TeachersHandler struct {
	db *sql.DB
}

func NewTeacherHandler(db *sql.DB) *TeachersHandler {
	return &TeachersHandler{db: db}
}

func (h *TeachersHandler) GetTeachersHandler(w http.ResponseWriter, _ *http.Request) {
	rows, err := h.db.Query("SELECT * FROM teachers;")
	if err != nil {
		http.Error(w, "Error retrieving all the teachers", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	teachersList := make([]models.Teacher, 0)
	for rows.Next() {
		var teacher models.Teacher
		err = rows.Scan(&teacher.Id, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
		if err != nil {
			http.Error(w, "Error scanning database results", http.StatusInternalServerError)
			return
		}
		teachersList = append(teachersList, teacher)
	}
	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(teachersList),
		Data:   teachersList,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *TeachersHandler) GetTeacherByIdhandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(err)
		return
	}

	var teacher models.Teacher
	err = h.db.QueryRow("SELECT * FROM teachers WHERE id = ?", id).Scan(
		&teacher.Id,
		&teacher.FirstName,
		&teacher.LastName,
		&teacher.Email,
		&teacher.Class,
		&teacher.Subject,
	)
	if err != nil {
		http.Error(w, "Teacher with that ID does not exist in the database!", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teacher)
}

func (h *TeachersHandler) CreateTeachersHandler(w http.ResponseWriter, r *http.Request) {
	var newTeachers []models.Teacher
	err := json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		http.Error(w, "Invalid Request body", http.StatusBadRequest)
		return
	}

	stmt, err := h.db.Prepare("INSERT INTO teachers(first_name, last_name, email, class, subject) VALUES(?,?,?,?,?)")
	if err != nil {
		http.Error(w, "Error in preparing SQL query", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	addedTeachers := make([]models.Teacher, len(newTeachers))
	for i, newTeacher := range newTeachers {
		res, err := stmt.Exec(newTeacher.FirstName, newTeacher.LastName, newTeacher.Email, newTeacher.Class, newTeacher.Subject)
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
		newTeacher.Id = int(lastId)
		addedTeachers[i] = newTeacher
	}

	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(addedTeachers),
		Data:   addedTeachers,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func (h *TeachersHandler) DeleteTeachersHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := h.db.Exec("DELETE FROM teachers WHERE id = ?", id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error deleting the teacher", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Println(err)
		http.Error(w, "Error deleting the teacher", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		log.Println(err)
		http.Error(w, "The teacher does not exist", http.StatusInternalServerError)
		return
	}

	response := struct {
		Status string `json:"status"`
		Id     int    `json:"id"`
	}{
		Status: "Successfully deleted the teacher",
		Id:     id,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *TeachersHandler) UpdateTeachersHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Teacher ID", http.StatusBadGateway)
		return
	}

	var updatedTeacher models.Teacher
	err = json.NewDecoder(r.Body).Decode(&updatedTeacher)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid request payload", http.StatusBadGateway)
		return
	}

	var existingTeacher models.Teacher
	err = h.db.QueryRow("SELECT * FROM teachers WHERE id = ?", id).Scan(
		&existingTeacher.Id,
		&existingTeacher.FirstName,
		&existingTeacher.LastName,
		&existingTeacher.Email,
		&existingTeacher.Class,
		&existingTeacher.Subject,
	)
	// In case there are no rows, an error is still thrown
	if err == sql.ErrNoRows {
		http.Error(w, "Teacher with the given ID not found!", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Unable to retrieve data", http.StatusInternalServerError)
		return
	}

	updatedTeacher.Id = existingTeacher.Id
	_, err = h.db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?",
		updatedTeacher.FirstName,
		updatedTeacher.LastName,
		updatedTeacher.Email,
		updatedTeacher.Class,
		updatedTeacher.Subject,
		id,
	)

	if err != nil {
		http.Error(w, "Error updating the teachers details", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTeacher)
}

func (h *TeachersHandler) GetStudentsByTeacherId(w http.ResponseWriter, r *http.Request) {
	// Allowed Roles: admin, manager, exec
	_, err := utils.AuthorizeUser(r.Context().Value(utils.ContextKey("role")).(string), "admin", "manager", "exec")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	teacherId := r.PathValue("id")
	var students []models.Student

	rows, err := h.db.Query("SELECT * FROM students WHERE class = (SELECT class FROM teachers WHERE id = ?)", teacherId)
	if err != nil {
		http.Error(w, "Error getting students under the given teacher", http.StatusNotFound)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var student models.Student
		err = rows.Scan(&student.Id, &student.FirstName, &student.LastName, &student.Email, &student.Class)
		if err != nil {
			http.Error(w, "Error querying the database", http.StatusInternalServerError)
			return
		}
		students = append(students, student)
	}

	err = rows.Err()
	if err != nil {
		http.Error(w, "Error querying the database", http.StatusInternalServerError)
		return
	}

	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Student `json:"data"`
	}{
		Status: "Success",
		Count:  len(students),
		Data:   students,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
