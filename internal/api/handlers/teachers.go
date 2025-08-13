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

func TeachersHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		getTeachersHandler(w, r)
		// path := strings.TrimPrefix(r.URL.Path, "/teachers/")
		// pathID := strings.TrimSuffix(path, "/")
		// fmt.Println("ID is: ", pathID)

		// queryParams := r.URL.Query()
		// sortBy := queryParams.Get("sortBy")
		// key := queryParams.Get("key")
		// sortorder := queryParams.Get("sortorder")

		// fmt.Printf("SortBy: %v, Key: %v, SortOrder: %v\n", sortBy, key, sortorder)

		w.Write([]byte("Hello get method on teachers route"))
		//fmt.Println("Hello get method on teachers route")
	case http.MethodPost:
		addTeachersHandler(w, r)
		w.Write([]byte("Hello post method on teachers route"))
		fmt.Println("Hello post method on teachers route")
	case http.MethodPut:
		w.Write([]byte("Hello put method on teachers route"))
		fmt.Println("Hello put method on teachers route")
	case http.MethodPatch:
		w.Write([]byte("Hello patch method on teachers route"))
		fmt.Println("Hello patch method on teachers route")
	case http.MethodDelete:
		w.Write([]byte("Hello delete method on teachers route"))
		fmt.Println("Hello delete method on teachers route")

	}
}

func getTeachersHandler(w http.ResponseWriter, r *http.Request) {

	db, err := sqlconnect.ConnectDb()
	if err != nil {
		http.Error(w, "Error connecting to db", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	query := "SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE 2=2"
	var args []interface{}

	params := map[string]string{
		"first_name": "first_name",
		"last_name":  "last_name",
		"email":      "email",
		"class":      "class",
		"subject":    "subject",
	}

	for param, dbField := range params {
		value := r.URL.Query().Get(param)
		//value := r.FormValue(param)
		if value != "" {
			query += " AND " + dbField + " = ?"
			args = append(args, value)
			fmt.Println("Query: ", query)
		}
		fmt.Println("Query: ", query)
		fmt.Println("Value: ", value)

	}
	// query += " AND " + "first_name" + " = ?"
	// args = append(args, "first_name")

	rows, err := db.Query(query, args...)
	if err != nil {
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	teachersList := make([]teacher.Teacher, 0)
	//teachersList := []teacher.Teacher{}

	for rows.Next() {
		var teacher teacher.Teacher
		err := rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
		if err != nil {
			http.Error(w, "Database scanning error", http.StatusInternalServerError)
			return
		}
		teachersList = append(teachersList, teacher)

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

func addTeachersHandler(w http.ResponseWriter, r *http.Request) {
	// mutex.Lock()
	// defer mutex.Unlock()

	db, err := sqlconnect.ConnectDb()
	if err != nil {
		http.Error(w, "Error connecting to db", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var newTeachers []teacher.Teacher
	err = json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		return
	}
	stmt, err := db.Prepare("INSERT INTO teachers (first_name, last_name, email, class, subject) VALUES(?,?,?,?,?)")
	if err != nil {
		http.Error(w, "Error to prepare db connection", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	addedTeachers := make([]teacher.Teacher, len(newTeachers))

	for i, newTeacher := range newTeachers {
		res, err := stmt.Exec(newTeacher.FirstName, newTeacher.LastName, newTeacher.Email, newTeacher.Class, newTeacher.Subject)
		if err != nil {
			http.Error(w, "Error to execute db operation", http.StatusInternalServerError)
			return
		}

		lastID, err := res.LastInsertId()
		if err != nil {
			http.Error(w, "Error for lastID", http.StatusInternalServerError)
			return
		}
		newTeacher.ID = int(lastID)
		addedTeachers[i] = newTeacher

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
