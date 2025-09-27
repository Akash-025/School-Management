package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	teacher "restapi/internal/models"
	"restapi/internal/repositories/sqlconnect"
	"strconv"
)

// var (
// 	teachers = make(map[int]teacher.Student)
// 	mutex    sync.Mutex
// 	next     = 1
// )

func GetStudentsHandler(w http.ResponseWriter, r *http.Request) {

	page, limit, teachersList, ok := sqlconnect.GetStudentsHandlerDB(w, r)
	if !ok {
		http.Error(w, "Error fo GetStudentsHandlerDB", http.StatusInternalServerError)
		return
	}
	respone := struct {
		Status   string            `json:"status"`
		Count    int               `json:"count"`
		Page     int               `json:"page"`
		PageSize int               `json:"page_size"`
		Data     []teacher.Student `json:"data"`
	}{
		Status:   "Success",
		Count:    len(teachersList),
		Page:     page,
		PageSize: limit,
		Data:     teachersList,
	}
	fmt.Println(respone)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respone)

}

func GetOneStudentHandler(w http.ResponseWriter, r *http.Request) {

	db, err := sqlconnect.ConnectDb()
	if err != nil {
		http.Error(w, "Error connecting to db", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return
	}

	var teacher teacher.Student
	err = db.QueryRow("SELECT id, first_name, last_name, email, class FROM students WHERE id = ?", id).Scan(
		&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class,
	)
	if err == sql.ErrNoRows {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database query Error ", http.StatusInternalServerError)
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(teacher)

}

func AddStudentsHandler(w http.ResponseWriter, r *http.Request) {
	// mutex.Lock()
	// defer mutex.Unlock()

	err, addedTeachers, ok := sqlconnect.AddStudentsHandlerDB(w, r)
	if !ok {
		http.Error(w, "Error for AddStudentsHandlerDB", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := struct {
		Status string            `json:"status"`
		Count  int               `json:"count"`
		Data   []teacher.Student `json:"data"`
	}{
		Status: "success",
		Count:  len(addedTeachers),
		Data:   addedTeachers,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}

}

func PatchStudentsHandler(w http.ResponseWriter, r *http.Request) {

	showUpdatedTeachers, ok := sqlconnect.PatchStudentsHandlerDB(w, r)
	if !ok {
		http.Error(w, "Error for PatchStudentsHandlerDB", http.StatusInternalServerError)
		return
	}

	respone := struct {
		Status string            `json:"status"`
		Count  int               `json:"count"`
		Data   []teacher.Student `json:"data"`
	}{
		Status: "Success",
		Count:  len(showUpdatedTeachers),
		Data:   showUpdatedTeachers,
	}

	json.NewEncoder(w).Encode(respone)
}

func DeleteStudentHandler(w http.ResponseWriter, r *http.Request) {

	db, err := sqlconnect.ConnectDb()
	if err != nil {
		http.Error(w, "Error connecting to db", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return
	}
	result, err := db.Exec("DELETE FROM students WHERE id=?", id)
	if err != nil {

		http.Error(w, "Error deleting from db", http.StatusInternalServerError)
		return
	}
	rowsAffect, err := result.RowsAffected()
	if err != nil {
		return
	}
	if rowsAffect == 0 {
		http.Error(w, "Student Not Found", http.StatusNotFound)
		return
	}

	response := struct {
		Status string `json:"status"`
	}{
		Status: "Successfully Deleted",
	}

	json.NewEncoder(w).Encode(response)
}

func GetStudentsByTeacherId(w http.ResponseWriter, r *http.Request) {

	students, ok := sqlconnect.GetStduntsByTeacherIdDB(r, w)
	if !ok {
		http.Error(w, "Error for GetStduntsByTeacherIdDB", http.StatusInternalServerError)
		return
	}

	response := struct {
		Status string
		Count  int
		Data   []teacher.Student
	}{
		Status: "Successful",
		Count:  len(students),
		Data:   students,
	}

	json.NewEncoder(w).Encode(response)
}
