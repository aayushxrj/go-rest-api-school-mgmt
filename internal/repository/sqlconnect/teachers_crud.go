package sqlconnect

import (
	"database/sql"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/aayushxrj/go-rest-api-school-mgmt/internal/models"
	"github.com/aayushxrj/go-rest-api-school-mgmt/pkg/utils"
)

func isValidSortOrder(order string) bool {
	return order == "asc" || order == "desc"
}
func isValidSortField(field string) bool {
	validFields := map[string]bool{
		"first_name": true,
		"last_name":  true,
		"email":      true,
		"class":      true,
		"subject":    true,
	}
	return validFields[field]
}

func addSorting(r *http.Request, query string) string {
	// teachers/?sortby=name:asc&sortby=class:desc
	sortParams := r.URL.Query()["sortby"]
	// fmt.Println(sortParams)

	if len(sortParams) > 0 {
		query += " ORDER BY"
		for i, param := range sortParams {
			parts := strings.Split(param, ":")
			if len(parts) != 2 {
				continue
			}
			field, order := parts[0], parts[1]
			if !isValidSortField(field) || !isValidSortOrder(order) {
				continue
			}
			if i > 0 {
				query += ","
			}
			query += " " + field + " " + order
		}
	}
	return query
}

func addFilters(r *http.Request, query string, args []any) (string, []any) {
	params := map[string]string{
		"first_name": "first_name",
		"last_name":  "last_name",
		"email":      "email",
		"class":      "class",
		"subject":    "subject",
	}

	for param, dbField := range params {
		value := r.URL.Query().Get(param)
		if value != "" {
			query += " AND " + dbField + " = ?"
			args = append(args, value)
		}
	}

	return query, args
}

func GetTeachersDBHandler(teachers []models.Teacher, r *http.Request) ([]models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Database connection error")
	}
	defer db.Close()

	//  Handle Query Parameters
	query := "SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE 1=1"
	var args []any

	query, args = addFilters(r, query, args)

	// teachers/?sortby=name:asc&sortby=class:desc
	query = addSorting(r, query)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, utils.ErrorHandler(err, "Database error")
	}
	defer rows.Close()

	// teacherList := []models.Teacher{}

	for rows.Next() {
		var teacher models.Teacher
		err := rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Database error")
		}
		teachers = append(teachers, teacher)
	}
	return teachers, nil
}

func GetOneTeacherDBHandler(id int) (models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "Database connection error")
	}
	defer db.Close()

	var teacher models.Teacher

	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(
		&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
	if err == sql.ErrNoRows {
		return models.Teacher{}, utils.ErrorHandler(err, "Teacher not found")
	} else if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "Database error")
	}
	return teacher, nil
}

func AddTeachersDBHandler(newTeachers []models.Teacher) ([]models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Database connection error")
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO teachers (first_name, last_name, email, class, subject) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return nil, utils.ErrorHandler(err, "Database error")
	}
	defer stmt.Close()

	addedTeachers := make([]models.Teacher, len(newTeachers))

	for i, newTeacher := range newTeachers {
		res, err := stmt.Exec(newTeacher.FirstName, newTeacher.LastName, newTeacher.Email, newTeacher.Class, newTeacher.Subject)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Database error")
		}
		lastID, err := res.LastInsertId()
		if err != nil {
			return nil, utils.ErrorHandler(err, "Database error")
		}
		newTeacher.ID = int(lastID)
		addedTeachers[i] = newTeacher
	}
	return addedTeachers, nil
}

func UpdateTeacherDBHandler(id int, updatedTeacher models.Teacher) (models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "Database connection error")
	}
	defer db.Close()

	var existingTeacher models.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(
		&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)
	if err == sql.ErrNoRows {
		return models.Teacher{}, utils.ErrorHandler(err, "Teacher not found")
	} else if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "Database error")
	}

	updatedTeacher.ID = existingTeacher.ID

	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?",
		updatedTeacher.FirstName, updatedTeacher.LastName, updatedTeacher.Email, updatedTeacher.Class, updatedTeacher.Subject, updatedTeacher.ID)
	if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "Database error")
	}
	return updatedTeacher, nil
}

func PatchTeachersDBHandler(updates []map[string]interface{}) error {
	db, err := ConnectDB()
	if err != nil {
		return utils.ErrorHandler(err, "Database connection error")
	}
	defer db.Close()

	// transaction
	tx, err := db.Begin()
	if err != nil {
		return utils.ErrorHandler(err, "Database error")
	}

	for _, update := range updates {
		// fmt.Println(update)
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

		var teacherFromDb models.Teacher
		err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(
			&teacherFromDb.ID, &teacherFromDb.FirstName, &teacherFromDb.LastName, &teacherFromDb.Email, &teacherFromDb.Class, &teacherFromDb.Subject)

		if err == sql.ErrNoRows {
			tx.Rollback()
			return utils.ErrorHandler(err, "Teacher not found with ID "+strconv.Itoa(id))
		} else if err != nil {
			tx.Rollback()
			return utils.ErrorHandler(err, "Database error")
		}

		// apply updates using reflect pkg
		teacherVal := reflect.ValueOf(&teacherFromDb).Elem()
		teacherType := teacherVal.Type()

		for k, v := range update {
			if k == "id" {
				continue
			}
			for i := 0; i < teacherType.NumField(); i++ {
				field := teacherType.Field(i)
				jsonTag := field.Tag.Get("json")
				if jsonTag == k+",omitempty" {
					fieldVal := teacherVal.Field(i)
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

		_, err = tx.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?",
			teacherFromDb.FirstName, teacherFromDb.LastName, teacherFromDb.Email, teacherFromDb.Class, teacherFromDb.Subject, teacherFromDb.ID)
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

func PatchOneTeacherDBHandler(id int, updates map[string]interface{}) (models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "Database connection error")
	}
	defer db.Close()

	var existingTeacher models.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(
		&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)
	if err == sql.ErrNoRows {
		return models.Teacher{}, utils.ErrorHandler(err, "Teacher not found")
	} else if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "Database error")
	}

	// apply updates
	// for k, v := range updates {
	// 	switch k {
	// 	case "first_name":
	// 		existingTeacher.FirstName = v.(string)
	// 	case "last_name":
	// 		existingTeacher.LastName = v.(string)
	// 	case "email":
	// 		existingTeacher.Email = v.(string)
	// 	case "class":
	// 		existingTeacher.Class = v.(string)
	// 	case "subject":
	// 		existingTeacher.Subject = v.(string)
	// 	}
	// }

	// Apply updates using reflect pkg
	teacherVal := reflect.ValueOf(&existingTeacher).Elem()
	// fmt.Println(teacherVal.Type())
	// fmt.Println(teacherVal.Type().Field(0))

	teacherType := teacherVal.Type()

	for k, v := range updates {
		for i := 0; i < teacherType.NumField(); i++ {
			field := teacherType.Field(i)
			jsonTag := field.Tag.Get("json")
			// fmt.Println(jsonTag)
			if jsonTag == k+",omitempty" {
				fieldVal := teacherVal.Field(i)
				if fieldVal.IsValid() && fieldVal.CanSet() {
					val := reflect.ValueOf(v)
					fieldVal.Set(val.Convert(fieldVal.Type()))
				}
				break
			}
		}
	}

	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?",
		existingTeacher.FirstName, existingTeacher.LastName, existingTeacher.Email, existingTeacher.Class, existingTeacher.Subject, existingTeacher.ID)
	if err != nil {
		return models.Teacher{}, utils.ErrorHandler(err, "Database error")
	}

	return existingTeacher, nil
}

func DeleteOneTeacherDBHandler(id int) error {
	db, err := ConnectDB()
	if err != nil {
		return utils.ErrorHandler(err, "Database connection error")
	}
	defer db.Close()

	res, err := db.Exec("DELETE FROM teachers WHERE id = ?", id)
	if err != nil {
		return utils.ErrorHandler(err, "Database error")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return utils.ErrorHandler(err, "Database error")
	}
	if rowsAffected == 0 {
		return utils.ErrorHandler(err, "Teacher not found")
	}
	return nil
}

func DeleteTeachersDBHandler(ids []int) ([]int, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Database connection error")
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Database error")
	}

	stmt, err := tx.Prepare("DELETE FROM teachers WHERE id = ?")
	if err != nil {
		tx.Rollback()
		return nil, utils.ErrorHandler(err, "Database error")
	}

	deletedIds := []int{}

	for _, id := range ids {
		res, err := stmt.Exec(id)
		if err != nil {
			tx.Rollback()
			return nil, utils.ErrorHandler(err, "Error deleteing teacher")
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			tx.Rollback()
			return nil, utils.ErrorHandler(err, "Database error")
		}
		if rowsAffected > 0 {
			deletedIds = append(deletedIds, id)
		}
		if rowsAffected < 1 {
			tx.Rollback()
			return nil, utils.ErrorHandler(err, "Teacher not found with ID "+strconv.Itoa(id))
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error committing transaction")
	}

	if len(deletedIds) < 1 {
		return nil, utils.ErrorHandler(err, "No teachers found to delete")
	}
	return deletedIds, nil
}
