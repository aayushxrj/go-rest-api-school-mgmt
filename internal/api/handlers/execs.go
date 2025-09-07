package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/aayushxrj/go-rest-api-school-mgmt/internal/models"
	"github.com/aayushxrj/go-rest-api-school-mgmt/internal/repository/sqlconnect"
)

// GetExecsHandler godoc
// @Summary Retrieve all execs
// @Description Get a list of execs with optional filtering and sorting.
// @Tags execs
// @Accept json
// @Produce json
// @Param first_name query string false "Filter by first name (optional)"
// @Param last_name query string false "Filter by last name (optional)"
// @Param email query string false "Filter by email (optional)"
// @Param role query string false "Filter by role (optional)"
// @Param sortby query string false "Sorting (e.g., first_name:asc, role:desc) (optional)"
// @Success 200 {object} map[string]interface{} "List of execs with metadata"
// @Failure 500 {string} string "Internal server error"
// @Router /execs [get]
func GetExecsHandler(w http.ResponseWriter, r *http.Request) {
	var execs []models.Exec
	execs, err := sqlconnect.GetExecsDBHandler(execs, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Status string         `json:"status"`
		Count  int            `json:"count"`
		Data   []models.Exec  `json:"data"`
	}{
		Status: "success",
		Count:  len(execs),
		Data:   execs,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetOneExecHandler godoc
// @Summary Get one exec
// @Description Retrieve details of an exec by ID
// @Tags execs
// @Produce json
// @Param id path int true "Exec ID"
// @Success 200 {object} models.Exec
// @Failure 400 {string} string "Invalid Exec ID"
// @Failure 500 {string} string "Internal server error"
// @Router /execs/{id} [get]
func GetOneExecHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Exec ID", http.StatusBadRequest)
		return
	}

	exec, err := sqlconnect.GetOneExecDBHandler(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(exec)
}

// AddExecHandler godoc
// @Summary Add new execs
// @Description Add one or more execs
// @Tags execs
// @Accept json
// @Produce json
// @Param execs body []models.Exec true "List of execs"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Internal server error"
// @Router /execs [post]
func AddExecHandler(w http.ResponseWriter, r *http.Request) {
	var rawExecs []map[string]interface{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body.", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &rawExecs)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	fields := GetFieldNames(models.Exec{})
	allowedFields := make(map[string]struct{})
	for _, field := range fields {
		allowedFields[field] = struct{}{}
	}

	for _, exec := range rawExecs {
		for key := range exec {
			if _, ok := allowedFields[key]; !ok {
				http.Error(w, "Unacceptable field found in request. Only use allowed fields", http.StatusBadRequest)
				return
			}
		}
	}

	var newExecs []models.Exec
	err = json.Unmarshal(body, &newExecs)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	for _, exec := range newExecs {
		if err := CheckBlankFields(exec); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	addedExecs, err := sqlconnect.AddExecsDBHandler(newExecs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string         `json:"status"`
		Count  int            `json:"count"`
		Data   []models.Exec  `json:"data"`
	}{
		Status: "success",
		Count:  len(addedExecs),
		Data:   addedExecs,
	}
	json.NewEncoder(w).Encode(response)
}

// PatchExecsHandler godoc
// @Summary Partially update multiple execs
// @Description Apply partial updates to multiple execs
// @Tags execs
// @Accept json
// @Param updates body []map[string]interface{} true "List of updates with exec IDs"
// @Success 204 "No Content"
// @Failure 400 {string} string "Invalid request payload"
// @Failure 500 {string} string "Internal server error"
// @Router /execs [patch]
func PatchExecsHandler(w http.ResponseWriter, r *http.Request) {
	var updates []map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = sqlconnect.PatchExecsDBHandler(updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// PatchOneExecHandler godoc
// @Summary Partially update one exec
// @Description Apply partial updates to a single exec by ID
// @Tags execs
// @Accept json
// @Produce json
// @Param id path int true "Exec ID"
// @Param updates body map[string]interface{} true "Partial updates"
// @Success 200 {object} models.Exec
// @Failure 400 {string} string "Invalid request payload or ID"
// @Failure 500 {string} string "Internal server error"
// @Router /execs/{id} [patch]
func PatchOneExecHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Exec ID", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	updatedExec, err := sqlconnect.PatchOneExecDBHandler(id, updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedExec)
}

// DeleteOneExecHandler godoc
// @Summary Delete one exec
// @Description Delete an exec by ID
// @Tags execs
// @Param id path int true "Exec ID"
// @Success 204 "No Content"
// @Failure 400 {string} string "Invalid Exec ID"
// @Failure 500 {string} string "Internal server error"
// @Router /execs/{id} [delete]
func DeleteOneExecHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Exec ID", http.StatusBadRequest)
		return
	}

	err = sqlconnect.DeleteOneExecDBHandler(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
