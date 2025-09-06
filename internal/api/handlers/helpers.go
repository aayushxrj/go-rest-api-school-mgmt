package handlers

import (
	"errors"
	"reflect"
	"strings"

	"github.com/aayushxrj/go-rest-api-school-mgmt/pkg/utils"
)

func CheckBlankFields(value interface{}) error {
	// if teacher.FirstName == "" || teacher.LastName == "" || teacher.Email == "" || teacher.Class == "" || teacher.Subject == "" {
	// 	http.Error(w, "All fields are required", http.StatusBadRequest)
	// 	return true
	// }

	val := reflect.ValueOf(value)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.Kind() == reflect.String && field.String() == "" {
			// fmt.Println("field.Kind():", field.Kind())
			// fmt.Println("reflect.String:", reflect.String)
			// fmt.Println("field.String():", field.String())
			// http.Error(w, "All fields are required", http.StatusBadRequest)
			return utils.ErrorHandler(errors.New("all fields are required"), "all fields are required")
		}
	}
	return nil
}

func GetFieldNames(model interface{}) []string {
	val := reflect.TypeOf(model)
	fields := []string{} // allowed fields

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldToAdd := strings.TrimSuffix(field.Tag.Get("json"), ",omitempty")
		fields = append(fields, fieldToAdd)
	}
	return fields
}
