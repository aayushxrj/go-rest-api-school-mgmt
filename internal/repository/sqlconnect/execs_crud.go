package sqlconnect

import (
	"database/sql"
	"net/http"
	"reflect"
	"strconv"

	"github.com/aayushxrj/go-rest-api-school-mgmt/internal/models"
	"github.com/aayushxrj/go-rest-api-school-mgmt/pkg/utils"
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
