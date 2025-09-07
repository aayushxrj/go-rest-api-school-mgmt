package utils

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

func GenerateInsertQuery(tableName string, model interface{}) string {
	modelType := reflect.TypeOf(model)
	var columns, placeholders string
	for i := 0; i < modelType.NumField(); i++ {
		dbTag := modelType.Field(i).Tag.Get("db")
		// fmt.Println("dbTag:", dbTag)
		dbTag = strings.TrimSuffix(dbTag, ",omitempty")
		if dbTag != "" && dbTag != "id" {
			if columns != "" {
				columns += ", "
				placeholders += ", "
			}
			columns += dbTag
			placeholders += "?"
		}
	}
	// fmt.Printf("INSERT INTO teachers (%s) VALUES (%s)", columns, placeholders)
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, columns, placeholders)
}

func GetStructValues(model interface{}) []interface{} {
	modelValue := reflect.ValueOf(model)
	modelType := reflect.TypeOf(model)
	values := []interface{}{}

	for i := 0; i < modelType.NumField(); i++ {
		dbTag := modelType.Field(i).Tag.Get("db")
		if dbTag != "" && dbTag != "id,omitempty" {
			values = append(values, modelValue.Field(i).Interface())
		}
	}
	// log.Printf("Values:", values)
	return values
}

func isValidSortOrder(order string) bool {
	return order == "asc" || order == "desc"
}

// func isValidSortField(field string) bool {
// 	validFields := map[string]bool{
// 		"first_name": true,
// 		"last_name":  true,
// 		"email":      true,
// 		"class":      true,
// 		"subject":    true,
// 	}
// 	return validFields[field]
// }

func isValidSortField(field string, model interface{}) bool {
	modelType := reflect.TypeOf(model)
	for i := 0; i < modelType.NumField(); i++ {
		dbTag := modelType.Field(i).Tag.Get("db")
		dbTag = strings.TrimSuffix(dbTag, ",omitempty")
		if dbTag == field {
			return true
		}
	}
	return false
}

func AddSorting(r *http.Request, query string, model interface{}) string {
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
			if !isValidSortField(field, model) || !isValidSortOrder(order) {
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

// func AddFilters(r *http.Request, query string, args []any) (string, []any) {
// 	params := map[string]string{
// 		"first_name": "first_name",
// 		"last_name":  "last_name",
// 		"email":      "email",
// 		"class":      "class",
// 		"subject":    "subject",
// 	}

// 	for param, dbField := range params {
// 		value := r.URL.Query().Get(param)
// 		if value != "" {
// 			query += " AND " + dbField + " = ?"
// 			args = append(args, value)
// 		}
// 	}

// 	return query, args
// }

func AddFilters(r *http.Request, query string, args []any, model interface{}) (string, []any) {
	modelType := reflect.TypeOf(model)

	for i := 0; i < modelType.NumField(); i++ {
		dbTag := modelType.Field(i).Tag.Get("db")
		dbTag = strings.TrimSuffix(dbTag, ",omitempty")
		if dbTag != "" && dbTag != "id" {
			value := r.URL.Query().Get(dbTag)
			if value != "" {
				query += " AND " + dbTag + " = ?"
				args = append(args, value)
			}
		}
	}

	return query, args
}
