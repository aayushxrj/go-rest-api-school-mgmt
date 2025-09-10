package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/aayushxrj/go-rest-api-school-mgmt/internal/models"
	"github.com/aayushxrj/go-rest-api-school-mgmt/internal/repository/sqlconnect"
	"github.com/aayushxrj/go-rest-api-school-mgmt/pkg/utils"
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
		Status string        `json:"status"`
		Count  int           `json:"count"`
		Data   []models.Exec `json:"data"`
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
		Status string        `json:"status"`
		Count  int           `json:"count"`
		Data   []models.Exec `json:"data"`
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

// LoginHandler godoc
// @Summary User Login
// @Description Authenticates an exec user using username and password and returns a JWT token.
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.Exec true "Login Credentials (username and password required)"
// @Success 200 {object} map[string]string "JWT token in response body and also set as HttpOnly cookie"
// @Failure 400 {string} string "Invalid request body or missing username/password"
// @Failure 403 {string} string "Account inactive or password incorrect"
// @Failure 500 {string} string "Could not create login token"
// @Router /execs/login [post]
func LoginHandler(w http.ResponseWriter, r *http.Request) {

	var req models.Exec

	// Data Validation
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password are blank", http.StatusBadRequest)
		return
	}

	// Search for user if user actually exists
	user, err := sqlconnect.LoginDBHandler(req.Username)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusBadRequest)
		return
	}

	// is user active
	if user.InactiveStatus {
		http.Error(w, "Account is inactive", http.StatusForbidden)
		return
	}

	// verify password
	err = utils.VerifyPassword(req.Password, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	// Generate JWT Token
	tokenString, err := utils.SignToken(user.ID, req.Username, user.Role)
	if err != nil {
		http.Error(w, "Could not create login token", http.StatusInternalServerError)
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
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "test",
		Value:    "testing",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(24 * time.Hour),
		SameSite: http.SameSiteStrictMode,
	})

	// Response Body
	w.Header().Set("Content-Type", "application/json")

	response := struct {
		Token string `json:"token"`
	}{
		Token: tokenString}
	json.NewEncoder(w).Encode(response)
}

// LogoutHandler godoc
// @Summary Log out a user
// @Description Logs out the currently authenticated user by clearing the JWT cookie.
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]string "Logged out successfully"
// @Failure 500 {string} string "Internal server error"
// @Router /execs/logout [post]
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
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
	w.Write([]byte(`{"message": "Logged out succesfully"}`))
}

// UpdatePasswordHandler godoc
// @Summary Update an exec's password
// @Description Allows an exec to update their password after providing the current password. A new JWT token is generated upon success.
// @Tags auth
// @Accept json
// @Produce json
// @Param id path int true "Exec ID"
// @Param body body models.UpdatePasswordRequest true "Password update request"
// @Success 200 {object} map[string]string "Password updated successfully"
// @Failure 400 {string} string "Invalid input or password update failed"
// @Failure 500 {string} string "Internal server error"
// @Router /execs/{id}/updatepassword [patch]
func UpdatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	userId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid exec ID", http.StatusBadRequest)
		return
	}

	var req models.UpdatePasswordRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}
	r.Body.Close()

	if req.CurrentPassword == "" || req.NewPassword == "" {
		http.Error(w, "Please enter password", http.StatusBadRequest)
		return
	}

	_, token, err := sqlconnect.UpdatePasswordDBHandler(userId, req.CurrentPassword, req.NewPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if token == "" {
		http.Error(w, "Password updated. Could not create token", http.StatusBadRequest)
		return
	}

	// Send token as a response or as a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "Bearer",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(24 * time.Hour),
		SameSite: http.SameSiteStrictMode,
	})

	// Response Body
	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Message string `json:"message"`
	}{
		Message: "Password updated successfully",
	}
	json.NewEncoder(w).Encode(response)
}

// ForgotPasswordHandler godoc
// @Summary Request password reset
// @Description Sends a password reset link to the exec's email.
// @Tags auth
// @Accept json
// @Produce plain
// @Param body body object{email=string} true "Exec email"
// @Success 200 {string} string "Password reset link sent"
// @Failure 400 {string} string "Invalid request or user not found"
// @Failure 500 {string} string "Internal server error"
// @Router /execs/forgotpassword [post]
func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}
	r.Body.Close()

	if req.Email == "" {
		http.Error(w, "Please enter the email", http.StatusBadRequest)
		return
	}

	err = sqlconnect.ForgotPasswordDBHandler(req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// respond with success confirmation
	fmt.Fprintf(w, "Password reset link sent to %s", req.Email)

}

// ResetPasswordHandler godoc
// @Summary Reset password using reset token
// @Description Resets the exec's password using a reset token sent via email.
// @Tags auth
// @Accept json
// @Produce plain
// @Param resetcode path string true "Password reset token"
// @Param body body object{new_password=string,confirm_password=string} true "New password request"
// @Success 200 {string} string "Password reset successfully"
// @Failure 400 {string} string "Invalid request or password mismatch"
// @Failure 500 {string} string "Internal server error"
// @Router /execs/resetpassword/reset/{resetcode} [post]
func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("resetcode")

	type request struct {
		NewPassword     string `json:"new_password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid values in request", http.StatusBadRequest)
		return
	}
	
	if req.ConfirmPassword == "" || req.NewPassword == "" {
		http.Error(w, "Please enter both the passwords", http.StatusBadRequest)
		return
	}

	if req.NewPassword != req.ConfirmPassword {
		http.Error(w, "Passwords should match", http.StatusBadRequest)
		return
	}

	// Hash the new password
	err = sqlconnect.ResetPasswordDBHandler(token, req.NewPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintln(w, "Password reset successfully")
}