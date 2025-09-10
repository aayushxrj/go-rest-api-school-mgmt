package sqlconnect

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/aayushxrj/go-rest-api-school-mgmt/internal/models"
	"github.com/aayushxrj/go-rest-api-school-mgmt/pkg/utils"
	"golang.org/x/crypto/argon2"

	"github.com/go-mail/mail/v2"
)

// GetExecsDBHandler retrieves a list of execs with optional filters and sorting
func GetExecsDBHandler(execs []models.Exec, r *http.Request) ([]models.Exec, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Database connection error")
	}
	defer db.Close()

	query := `SELECT id, first_name, last_name, email, username, user_created_at, inactive_status, role FROM execs WHERE 1=1`
	var args []any

	query, args = utils.AddFilters(r, query, args, models.Exec{})
	query = utils.AddSorting(r, query, models.Exec{})

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Database error")
	}
	defer rows.Close()

	for rows.Next() {
		var exec models.Exec
		err := rows.Scan(
			&exec.ID, &exec.FirstName, &exec.LastName, &exec.Email,
			&exec.Username, &exec.UserCreatedAt, &exec.InactiveStatus, &exec.Role,
		)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Database error")
		}
		execs = append(execs, exec)
	}
	return execs, nil
}

// GetOneExecDBHandler retrieves a single exec by ID
func GetOneExecDBHandler(id int) (models.Exec, error) {
	db, err := ConnectDB()
	if err != nil {
		return models.Exec{}, utils.ErrorHandler(err, "Database connection error")
	}
	defer db.Close()

	var exec models.Exec
	err = db.QueryRow(`SELECT id, first_name, last_name, email, username, user_created_at, inactive_status, role FROM execs WHERE id = ?`, id).Scan(
		&exec.ID, &exec.FirstName, &exec.LastName, &exec.Email,
		&exec.Username, &exec.UserCreatedAt, &exec.InactiveStatus, &exec.Role,
	)
	if err == sql.ErrNoRows {
		return models.Exec{}, utils.ErrorHandler(err, "Exec not found")
	} else if err != nil {
		return models.Exec{}, utils.ErrorHandler(err, "Database error")
	}
	return exec, nil
}

// AddExecsDBHandler inserts new execs
func AddExecsDBHandler(newExecs []models.Exec) ([]models.Exec, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Database connection error")
	}
	defer db.Close()

	stmt, err := db.Prepare(utils.GenerateInsertQuery("execs", models.Exec{}))
	if err != nil {
		return nil, utils.ErrorHandler(err, "Database error")
	}
	defer stmt.Close()

	addedExecs := make([]models.Exec, len(newExecs))

	for i, newExec := range newExecs {

		if newExec.Password == "" {
			return nil, utils.ErrorHandler(errors.New("password is blank"), "Please enter the password")
		}
		salt := make([]byte, 16)
		_, err := rand.Read(salt)
		if err != nil {
			return nil, utils.ErrorHandler(errors.New("failed to generate salt"), "Database error")
		}

		hash := argon2.IDKey([]byte(newExec.Password), salt, 1, 64*1024, 4, 32)
		saltBase64 := base64.StdEncoding.EncodeToString(salt)
		hashBase64 := base64.StdEncoding.EncodeToString(hash)
		encodedHash := fmt.Sprintf("%s.%s", saltBase64, hashBase64)

		newExec.Password = encodedHash

		values := utils.GetStructValues(newExec)
		res, err := stmt.Exec(values...)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Database error")
		}
		lastID, err := res.LastInsertId()
		if err != nil {
			return nil, utils.ErrorHandler(err, "Database error")
		}
		newExec.ID = int(lastID)
		addedExecs[i] = newExec
	}
	return addedExecs, nil
}

// PatchExecsDBHandler performs partial updates for multiple execs
func PatchExecsDBHandler(updates []map[string]interface{}) error {
	db, err := ConnectDB()
	if err != nil {
		return utils.ErrorHandler(err, "Database connection error")
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return utils.ErrorHandler(err, "Database error")
	}

	for _, update := range updates {
		idStr, ok := update["id"].(string)
		if !ok {
			tx.Rollback()
			return utils.ErrorHandler(err, "Invalid or missing ID in update object")
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			tx.Rollback()
			return utils.ErrorHandler(err, "Error converting string to int")
		}

		var execFromDb models.Exec
		err = db.QueryRow(`SELECT id, first_name, last_name, email, username FROM execs WHERE id = ?`, id).Scan(
			&execFromDb.ID, &execFromDb.FirstName, &execFromDb.LastName, &execFromDb.Email, &execFromDb.Username,
		)
		if err == sql.ErrNoRows {
			tx.Rollback()
			return utils.ErrorHandler(err, "Exec not found with ID "+strconv.Itoa(id))
		} else if err != nil {
			tx.Rollback()
			return utils.ErrorHandler(err, "Database error")
		}

		execVal := reflect.ValueOf(&execFromDb).Elem()
		execType := execVal.Type()

		for k, v := range update {
			if k == "id" {
				continue
			}
			for i := 0; i < execType.NumField(); i++ {
				field := execType.Field(i)
				jsonTag := field.Tag.Get("json")
				if jsonTag == k+",omitempty" {
					fieldVal := execVal.Field(i)
					if fieldVal.IsValid() && fieldVal.CanSet() {
						val := reflect.ValueOf(v)
						if val.Type().ConvertibleTo(fieldVal.Type()) {
							fieldVal.Set(val.Convert(fieldVal.Type()))
						} else {
							tx.Rollback()
							return utils.ErrorHandler(err, "Type mismatch for field "+k)
						}
					}
					break
				}
			}
		}

		_, err = tx.Exec(`UPDATE execs SET first_name=?, last_name=?, email=?, username=? WHERE id=?`,
			execFromDb.FirstName, execFromDb.LastName, execFromDb.Email,
			execFromDb.Username, execFromDb.ID)
		if err != nil {
			tx.Rollback()
			return utils.ErrorHandler(err, "Database error")
		}
	}

	err = tx.Commit()
	if err != nil {
		return utils.ErrorHandler(err, "Error committing transaction")
	}
	return nil
}

// PatchOneExecDBHandler performs partial update for one exec
func PatchOneExecDBHandler(id int, updates map[string]interface{}) (models.Exec, error) {
	db, err := ConnectDB()
	if err != nil {
		return models.Exec{}, utils.ErrorHandler(err, "Database connection error")
	}
	defer db.Close()

	var existingExec models.Exec
	err = db.QueryRow(`SELECT id, first_name, last_name, email, username FROM execs WHERE id = ?`, id).Scan(
		&existingExec.ID, &existingExec.FirstName, &existingExec.LastName, &existingExec.Email, &existingExec.Username,
	)
	if err == sql.ErrNoRows {
		return models.Exec{}, utils.ErrorHandler(err, "Exec not found")
	} else if err != nil {
		return models.Exec{}, utils.ErrorHandler(err, "Database error")
	}

	execVal := reflect.ValueOf(&existingExec).Elem()
	execType := execVal.Type()

	for k, v := range updates {
		for i := 0; i < execType.NumField(); i++ {
			field := execType.Field(i)
			jsonTag := field.Tag.Get("json")
			if jsonTag == k+",omitempty" {
				fieldVal := execVal.Field(i)
				if fieldVal.IsValid() && fieldVal.CanSet() {
					val := reflect.ValueOf(v)
					fieldVal.Set(val.Convert(fieldVal.Type()))
				}
				break
			}
		}
	}

	_, err = db.Exec(`UPDATE execs 
		SET first_name=?, last_name=?, email=?, username=? WHERE id=?`,
		existingExec.FirstName, existingExec.LastName, existingExec.Email, existingExec.Username, existingExec.ID)
	if err != nil {
		return models.Exec{}, utils.ErrorHandler(err, "Database error")
	}

	return existingExec, nil
}

// DeleteOneExecDBHandler deletes a single exec
func DeleteOneExecDBHandler(id int) error {
	db, err := ConnectDB()
	if err != nil {
		return utils.ErrorHandler(err, "Database connection error")
	}
	defer db.Close()

	res, err := db.Exec("DELETE FROM execs WHERE id = ?", id)
	if err != nil {
		return utils.ErrorHandler(err, "Database error")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return utils.ErrorHandler(err, "Database error")
	}
	if rowsAffected == 0 {
		return utils.ErrorHandler(err, "Exec not found")
	}
	return nil
}

func LoginDBHandler(username string) (*models.Exec, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Database connection error")
	}
	defer db.Close()

	user := &models.Exec{}
	err = db.QueryRow(`SELECT id, first_name, last_name, email, username, password, inactive_status, role FROM execs WHERE username = ?`, username).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Email,
		&user.Username, &user.Password, &user.InactiveStatus, &user.Role,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, utils.ErrorHandler(err, "User not found")
		}
		return nil, utils.ErrorHandler(err, "Database error")
	}
	return user, nil
}

func UpdatePasswordDBHandler(userId int, currentPassword, newPassword string) (bool, string, error) {
	db, err := ConnectDB()
	if err != nil {
		return false, "", utils.ErrorHandler(err, "database connection error")
	}
	defer db.Close()

	var username string
	var userPassword string
	var userRole string

	err = db.QueryRow("SELECT username, password, role FROM execs WHERE id = ?", userId).Scan(&username, &userPassword, &userRole)
	if err != nil {
		return false, "", utils.ErrorHandler(err, "user not found")
	}

	err = utils.VerifyPassword(currentPassword, userPassword)
	if err != nil {
		return false, "", utils.ErrorHandler(err, "The password you entered does not match the current password on file.")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return false, "", utils.ErrorHandler(err, "internal error")
	}

	currentTime := time.Now().Format(time.RFC3339)

	_, err = db.Exec("UPDATE execs SET password = ?, password_changed_at = ? WHERE id = ?", hashedPassword, currentTime, userId)
	if err != nil {
		return false, "", utils.ErrorHandler(err, "failed to update the password")
	}

	token, err := utils.SignToken(userId, username, userRole)
	if err != nil {
		utils.ErrorHandler(err, "Password updated. Could not create token")
		return false, "", utils.ErrorHandler(err, "Password updated. Could not create token")
	}

	return true, token, nil
}

func ForgotPasswordDBHandler(emailId string) error {
	db, err := ConnectDB()
	if err != nil {
		return utils.ErrorHandler(err, "Internal error")
	}
	defer db.Close()

	var exec models.Exec
	err = db.QueryRow("SELECT id FROM execs WHERE email = ?", emailId).Scan(&exec.ID)
	if err != nil {
		return utils.ErrorHandler(err, "User not found")
	}

	duration, err := strconv.Atoi(os.Getenv("RESET_TOKEN_EXP_DURATION"))
	if err != nil {
		return utils.ErrorHandler(err, "Failed to send password reset email")
	}
	mins := time.Duration(duration)

	expiry := time.Now().Add(mins * time.Minute).Format(time.RFC3339)

	tokenBytes := make([]byte, 32)
	_, err = rand.Read(tokenBytes)
	if err != nil {
		return utils.ErrorHandler(err, "Failed to send password reset email")
	}

	token := hex.EncodeToString(tokenBytes)

	hashedToken := sha256.Sum256(tokenBytes)

	hashedTokenString := hex.EncodeToString(hashedToken[:])

	_, err = db.Exec("UPDATE execs SET password_reset_token = ?, password_token_expires = ? WHERE id = ?", hashedTokenString, expiry, exec.ID)
	if err != nil {
		return utils.ErrorHandler(err, "Failed to send password reset email")
	}

	resetURL := fmt.Sprintf("https://localhost:3000/execs/resetpassword/reset/%s", token)
	message := fmt.Sprintf("Forgot your password?  Reset your password using the following link: \n%s\nIf you didn't request a password reset, please ignore this email. This link is only valid for %d minutes.", resetURL, int(mins))

	m := mail.NewMessage()
	m.SetHeader("From", "schooladmin@shool.com")
	m.SetHeader("To", emailId)
	m.SetHeader("Subject", "Your password reset link")
	m.SetBody("text/plain", message)

	d := mail.NewDialer("localhost", 1025, "", "")
	err = d.DialAndSend(m)
	if err != nil {
		return utils.ErrorHandler(err, "Failed to send password reset email")
	}
	
	return nil
}
