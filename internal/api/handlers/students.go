package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/aayushxrj/go-rest-api-school-mgmt/internal/models"
	"github.com/aayushxrj/go-rest-api-school-mgmt/internal/repository/sqlconnect"
)

// GetStudentsHandler godoc
// @Summary Retrieve all students
// @Description Get a list of students with optional filtering and sorting.
// @Tags students
// @Accept json
// @Produce json
// @Param first_name query string false "Filter by first name (optional)"
// @Param last_name query string false "Filter by last name (optional)"
// @Param email query string false "Filter by email (optional)"
// @Param class query string false "Filter by class (optional)"
// @Param sortby query string false "Sorting (e.g., first_name:asc, class:desc) (optional)"
// @Success 200 {object} map[string]interface{} "List of students with metadata"
// @Failure 500 {string} string "Internal server error"
// @Router /students [get]
func GetStudentsHandler(w http.ResponseWriter, r *http.Request) {
	var students []models.Student
	students, err := sqlconnect.GetStudentsDBHandler(students, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Student `json:"data"`
	}{
		Status: "success",
		Count:  len(students),
		Data:   students,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetOneStudentHandler godoc
// @Summary Get one student
// @Description Retrieve details of a student by ID
// @Tags students
// @Produce json
// @Param id path int true "Student ID"
// @Success 200 {object} models.Student
// @Failure 400 {string} string "Invalid Student ID"
// @Failure 500 {string} string "Internal server error"
// @Router /students/{id} [get]
func GetOneStudentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Student ID", http.StatusBadRequest)
		return
	}

	student, err := sqlconnect.GetOneStudentDBHandler(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(student)
}

// AddStudentHandler godoc
// @Summary Add new students
// @Description Add one or more students
// @Tags students
// @Accept json
// @Produce json
// @Param students body []models.Student true "List of students"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Internal server error"
// @Router /students [post]
func AddStudentHandler(w http.ResponseWriter, r *http.Request) {
	var rawStudents []map[string]interface{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body.", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &rawStudents)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	fields := GetFieldNames(models.Student{})
	allowedFields := make(map[string]struct{})
	for _, field := range fields {
		allowedFields[field] = struct{}{}
	}

	for _, student := range rawStudents {
		for key := range student {
			_, ok := allowedFields[key]
			if !ok {
				http.Error(w, "Unacceptable field found in request. Only use allowed fields", http.StatusBadRequest)
				return
			}
		}
	}

	var newStudents []models.Student
	err = json.Unmarshal(body, &newStudents)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	for _, student := range newStudents {
		err := CheckBlankFields(student)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	addedStudents, err := sqlconnect.AddStudentsDBHandler(newStudents)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Student `json:"data"`
	}{
		Status: "success",
		Count:  len(addedStudents),
		Data:   addedStudents,
	}
	json.NewEncoder(w).Encode(response)
}

// UpdateStudentHandler godoc
// @Summary Update a student
// @Description Update an existing student by ID
// @Tags students
// @Accept json
// @Produce json
// @Param id path int true "Student ID"
// @Param student body models.Student true "Updated student"
// @Success 200 {object} models.Student
// @Failure 400 {string} string "Invalid request payload or ID"
// @Failure 500 {string} string "Internal server error"
// @Router /students/{id} [put]
func UpdateStudentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Student ID", http.StatusBadRequest)
		return
	}

	var updatedStudent models.Student
	err = json.NewDecoder(r.Body).Decode(&updatedStudent)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	updatedStudentFromDB, err := sqlconnect.UpdateStudentDBHandler(id, updatedStudent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedStudentFromDB)
}

// PatchStudentsHandler godoc
// @Summary Partially update multiple students
// @Description Apply partial updates to multiple students
// @Tags students
// @Accept json
// @Param updates body []map[string]interface{} true "List of updates"
// @Success 204 "No Content"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Internal server error"
// @Router /students [patch]
func PatchStudentsHandler(w http.ResponseWriter, r *http.Request) {
	var updates []map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = sqlconnect.PatchStudentsDBHandler(updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// PatchOneStudentHandler godoc
// @Summary Partially update one student
// @Description Apply partial updates to a single student by ID
// @Tags students
// @Accept json
// @Produce json
// @Param id path int true "Student ID"
// @Param updates body map[string]interface{} true "Partial updates"
// @Success 200 {object} models.Student
// @Failure 400 {string} string "Invalid request payload or ID"
// @Failure 500 {string} string "Internal server error"
// @Router /students/{id} [patch]
func PatchOneStudentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Student ID", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	updatedStudent, err := sqlconnect.PatchOneStudentDBHandler(id, updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedStudent)
}

// DeleteOneStudentHandler godoc
// @Summary Delete one student
// @Description Delete a student by ID
// @Tags students
// @Param id path int true "Student ID"
// @Success 204 "No Content"
// @Failure 400 {string} string "Invalid Student ID"
// @Failure 500 {string} string "Internal server error"
// @Router /students/{id} [delete]
func DeleteOneStudentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Student ID", http.StatusBadRequest)
		return
	}

	err = sqlconnect.DeleteOneStudentDBHandler(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteStudentsHandler godoc
// @Summary Delete multiple students
// @Description Delete multiple students by their IDs
// @Tags students
// @Accept json
// @Produce json
// @Param ids body []int true "List of student IDs"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Internal server error"
// @Router /students [delete]
func DeleteStudentsHandler(w http.ResponseWriter, r *http.Request) {
	var ids []int
	err := json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	deletedIds, err := sqlconnect.DeleteStudentsDBHandler(ids)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := struct {
		Status     string `json:"status"`
		DeletedIDs []int  `json:"deleted_ids"`
	}{
		Status:     "Students deleted successfully",
		DeletedIDs: deletedIds,
	}
	json.NewEncoder(w).Encode(response)
}
