package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/aayushxrj/go-rest-api-school-mgmt/internal/models"
	"github.com/aayushxrj/go-rest-api-school-mgmt/internal/repository/sqlconnect"
)

var (
	teachers = make(map[int]models.Teacher)
	mutex    = &sync.Mutex{} // use in the post method
	nextID   = 1
)

func init() {
	teachers[nextID] = models.Teacher{ID: nextID, FirstName: "John", LastName: "Doe", Class: "10", Subject: "Math"}
	nextID++
	teachers[nextID] = models.Teacher{ID: nextID, FirstName: "Jane", LastName: "Smith", Class: "9", Subject: "Science"}
	nextID++
	teachers[nextID] = models.Teacher{ID: nextID, FirstName: "Jane", LastName: "Doe", Class: "11", Subject: "Biology"}
	nextID++
}

func TeachersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTeachersHandler(w, r)
	case http.MethodPost:
		addTeacherHandler(w, r)
	case http.MethodPut:
		w.Write([]byte("PUT Method on Teachers Route"))
	case http.MethodPatch:
		w.Write([]byte("PATCH Method on Teachers Route"))
	case http.MethodDelete:
		w.Write([]byte("DELETE Method on Teachers Route"))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed"))
	}
}

func getTeachersHandler(w http.ResponseWriter, r *http.Request) {

	// Path Parameters
	pathStr := strings.TrimPrefix(r.URL.Path, "/teachers/")
	idStr := strings.TrimSuffix(pathStr, "/")
	// fmt.Println("ID String:", idStr)

	// Query Parameters
	if idStr == "" {
		first_name := r.URL.Query().Get("first_name")
		last_name := r.URL.Query().Get("last_name")

		teacherList := make([]models.Teacher, 0, len(teachers))
		for _, teacher := range teachers {
			if (first_name == "" || first_name == teacher.FirstName) && (last_name == "" || last_name == teacher.LastName) {
				teacherList = append(teacherList, teacher)
			}
		}

		response := struct {
			Status string           `json:"status"`
			Count  int              `json:"count"`
			Data   []models.Teacher `json:"data"`
		}{
			Status: "success",
			Count:  len(teacherList),
			Data:   teacherList,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Handle Path Parameter
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid models.Teacher ID", http.StatusBadRequest)
		return
	}

	teacher, exists := teachers[id]
	if !exists {
		http.Error(w, "models.Teacher Not Found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teacher)
}

func addTeacherHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var newTeachers []models.Teacher
	err = json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("INSERT INTO teachers (first_name, last_name, email, class, subject) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	addedTeachers := make([]models.Teacher, len(newTeachers))

	for i, newTeacher := range newTeachers {
		res, err := stmt.Exec(newTeacher.FirstName, newTeacher.LastName, newTeacher.Email, newTeacher.Class, newTeacher.Subject)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		lastID, err := res.LastInsertId()
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		newTeacher.ID = int(lastID)
		addedTeachers[i] = newTeacher
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
