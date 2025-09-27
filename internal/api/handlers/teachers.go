package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	teacher "restapi/internal/models"
	"restapi/internal/repositories/sqlconnect"
	"strconv"
	"sync"
)

var (
	teachers = make(map[int]teacher.Teacher)
	mutex    sync.Mutex
	next     = 1
)

func GetTeachersHandler(w http.ResponseWriter, r *http.Request) {

	teachersList, ok := sqlconnect.GetTeachersHandlerDB(w, r)
	if !ok {
		http.Error(w, "Error for GetTeachersHandlerDB", http.StatusInternalServerError)
		return
	}
	respone := struct {
		Status string            `json:"status"`
		Count  int               `json:"count"`
		Data   []teacher.Teacher `json:"data"`
	}{
		Status: "Success",
		Count:  len(teachersList),
		Data:   teachersList,
	}
	fmt.Println(respone)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respone)

}

func GetOneTeacherHandler(w http.ResponseWriter, r *http.Request) {

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

	var teacher teacher.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(
		&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject,
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

func AddTeachersHandler(w http.ResponseWriter, r *http.Request) {
	// mutex.Lock()
	// defer mutex.Unlock()

	err, addedTeachers, ok := sqlconnect.AddTeachersHandlerDB(w, r)
	if !ok {
		http.Error(w, "Error for AddTeachersHandlerDB", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := struct {
		Status string            `json:"status"`
		Count  int               `json:"count"`
		Data   []teacher.Teacher `json:"data"`
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

func PatchTeachersHandler(w http.ResponseWriter, r *http.Request) {

	showUpdatedTeachers, ok := sqlconnect.PatchTeachersHandlerDB(w, r)
	if !ok {
		http.Error(w, "Error for PatchTeachersHandlerDB", http.StatusInternalServerError)
		return
	}

	respone := struct {
		Status string            `json:"status"`
		Count  int               `json:"count"`
		Data   []teacher.Teacher `json:"data"`
	}{
		Status: "Success",
		Count:  len(showUpdatedTeachers),
		Data:   showUpdatedTeachers,
	}

	json.NewEncoder(w).Encode(respone)
}

func DeleteTeacherHandler(w http.ResponseWriter, r *http.Request) {

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
	result, err := db.Exec("DELETE FROM teachers WHERE id=?", id)
	if err != nil {

		http.Error(w, "Error deleting from db", http.StatusInternalServerError)
		return
	}
	rowsAffect, err := result.RowsAffected()
	if err != nil {
		return
	}
	if rowsAffect == 0 {
		http.Error(w, "Teacher Not Found", http.StatusNotFound)
		return
	}

	response := struct {
		Status string `json:"status"`
	}{
		Status: "Successfully Deleted",
	}

	json.NewEncoder(w).Encode(response)
}
