package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/aayushxrj/go-rest-api-school-mgmt/internal/models"
	"github.com/aayushxrj/go-rest-api-school-mgmt/internal/repository/sqlconnect"
)

// GetTeachersHandler godoc
// @Summary Retrieve all teachers
// @Description Get a list of teachers with optional filtering and sorting.
// @Tags teachers
// @Accept json
// @Produce json
// @Param first_name query string false "Filter by first name (optional)"
// @Param last_name query string false "Filter by last name (optional)"
// @Param email query string false "Filter by email (optional)"
// @Param class query string false "Filter by class (optional)"
// @Param subject query string false "Filter by subject (optional)"
// @Param sortby query string false "Sorting (e.g., first_name:asc, class:desc) (optional)"
// @Success 200 {object} map[string]interface{} "List of teachers with metadata"
// @Failure 500 {string} string "Internal server error"
// @Router /teachers [get]
func GetTeachersHandler(w http.ResponseWriter, r *http.Request) {
	var teachers []models.Teacher
	teachers, err := sqlconnect.GetTeachersDBHandler(teachers, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(teachers),
		Data:   teachers,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetOneTeacherHandler godoc
// @Summary Get one teacher
// @Description Retrieve details of a teacher by ID
// @Tags teachers
// @Produce json
// @Param id path int true "Teacher ID"
// @Success 200 {object} models.Teacher
// @Failure 400 {string} string "Invalid Teacher ID"
// @Failure 500 {string} string "Internal server error"
// @Router /teachers/{id} [get]
func GetOneTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Teacher ID", http.StatusBadRequest)
		return
	}

	teacher, err := sqlconnect.GetOneTeacherDBHandler(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(teacher)
}

// AddTeacherHandler godoc
// @Summary Add new teachers
// @Description Add one or more teachers
// @Tags teachers
// @Accept json
// @Produce json
// @Param teachers body []models.Teacher true "List of teachers"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Internal server error"
// @Router /teachers [post]
func AddTeacherHandler(w http.ResponseWriter, r *http.Request) {

	// validate the data before sending furthur
	var rawTeachers []map[string]interface{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body.", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &rawTeachers)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// for non-existent fields
	fields := GetFieldNames(models.Teacher{})

	allowedFields := make(map[string]struct{})
	for _, field := range fields {
		allowedFields[field] = struct{}{}
	}

	for _, teacher := range rawTeachers {
		for key := range teacher {
			_, ok := allowedFields[key]
			if !ok {
				http.Error(w, "Unacceptable field found in request. Only use allowed fields", http.StatusBadRequest)
				return
			}
		}
	}

	var newTeachers []models.Teacher
	err = json.Unmarshal(body, &newTeachers)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// for blank value in fields
	for _, teacher := range newTeachers {
		err := CheckBlankFields(teacher)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	addedTeachers, err := sqlconnect.AddTeachersDBHandler(newTeachers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(addedTeachers),
		Data:   addedTeachers,
	}
	json.NewEncoder(w).Encode(response)
}

// UpdateTeacherHandler godoc
// @Summary Update a teacher
// @Description Update an existing teacher by ID
// @Tags teachers
// @Accept json
// @Produce json
// @Param id path int true "Teacher ID"
// @Param teacher body models.Teacher true "Updated teacher"
// @Success 200 {object} models.Teacher
// @Failure 400 {string} string "Invalid request payload or ID"
// @Failure 500 {string} string "Internal server error"
// @Router /teachers/{id} [put]
func UpdateTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Teacher ID", http.StatusBadRequest)
		return
	}

	var updatedTeacher models.Teacher
	err = json.NewDecoder(r.Body).Decode(&updatedTeacher)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	updatedTeacherFromDB, err := sqlconnect.UpdateTeacherDBHandler(id, updatedTeacher)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedTeacherFromDB)
}

// PatchTeachersHandler godoc
// @Summary Partially update multiple teachers
// @Description Apply partial updates to multiple teachers
// @Tags teachers
// @Accept json
// @Param updates body []map[string]interface{} true "List of updates"
// @Success 204 "No Content"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Internal server error"
// @Router /teachers [patch]
func PatchTeachersHandler(w http.ResponseWriter, r *http.Request) {
	var updates []map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = sqlconnect.PatchTeachersDBHandler(updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// PatchOneTeacherHandler godoc
// @Summary Partially update one teacher
// @Description Apply partial updates to a single teacher by ID
// @Tags teachers
// @Accept json
// @Produce json
// @Param id path int true "Teacher ID"
// @Param updates body map[string]interface{} true "Partial updates"
// @Success 200 {object} models.Teacher
// @Failure 400 {string} string "Invalid request payload or ID"
// @Failure 500 {string} string "Internal server error"
// @Router /teachers/{id} [patch]
func PatchOneTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Teacher ID", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	updatedTeacher, err := sqlconnect.PatchOneTeacherDBHandler(id, updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedTeacher)
}

// DeleteOneTeacherHandler godoc
// @Summary Delete one teacher
// @Description Delete a teacher by ID
// @Tags teachers
// @Param id path int true "Teacher ID"
// @Success 204 "No Content"
// @Failure 400 {string} string "Invalid Teacher ID"
// @Failure 500 {string} string "Internal server error"
// @Router /teachers/{id} [delete]
func DeleteOneTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Teacher ID", http.StatusBadRequest)
		return
	}

	err = sqlconnect.DeleteOneTeacherDBHandler(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteTeachersHandler godoc
// @Summary Delete multiple teachers
// @Description Delete multiple teachers by their IDs
// @Tags teachers
// @Accept json
// @Produce json
// @Param ids body []int true "List of teacher IDs"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Internal server error"
// @Router /teachers [delete]
func DeleteTeachersHandler(w http.ResponseWriter, r *http.Request) {
	var ids []int
	err := json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	deletedIds, err := sqlconnect.DeleteTeachersDBHandler(ids)
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
		Status:     "Teachers deleted successfully",
		DeletedIDs: deletedIds,
	}
	json.NewEncoder(w).Encode(response)
}

// GetStudentsByTeacherIDHandler godoc
// @Summary Retrieve students by teacher ID
// @Description Get all students assigned to a specific teacher.
// @Tags teachers
// @Accept json
// @Produce json
// @Param id path int true "Teacher ID"
// @Success 200 {object} map[string]interface{} "List of students with metadata"
// @Failure 400 {string} string "Invalid Teacher ID"
// @Failure 500 {string} string "Internal server error"
// @Router /teachers/{id}/students [get]
func GetStudentsByTeacherIDHandler(w http.ResponseWriter, r *http.Request){
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Teacher ID", http.StatusBadRequest)
		return
	}

	var students []models.Student

	students, err = sqlconnect.GetStudentsByTeacherIdDBHandler(id, students)
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

// GetStudentsCountByTeacherIDHandler godoc
// @Summary Retrieve student count by teacher ID
// @Description Get the total number of students assigned to a specific teacher.
// @Tags teachers
// @Accept json
// @Produce json
// @Param id path int true "Teacher ID"
// @Success 200 {object} map[string]interface{} "Student count"
// @Failure 400 {string} string "Invalid Teacher ID"
// @Failure 500 {string} string "Internal server error"
// @Router /teachers/{id}/studentcount [get]
func GetStudentsCountByTeacherIDHandler(w http.ResponseWriter, r *http.Request) {

	teacherId := r.PathValue("id")

	studentCount, err := sqlconnect.GetStudentsCountByTeacherIdDBHandler(teacherId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Status string `json:"status"`
		Count  int    `json:"count"`
	}{
		Status: "success",
		Count:  studentCount,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}