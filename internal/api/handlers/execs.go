package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"restapi/internal/models"
	teacher "restapi/internal/models"
	"restapi/internal/repositories/sqlconnect"
	"strconv"
	"time"
)

func GetExecsHandler(w http.ResponseWriter, r *http.Request) {

	teachersList, ok := sqlconnect.GetExecsHandlerDB(w, r)
	if !ok {
		http.Error(w, "Error for GetExecsHandlerDB", http.StatusInternalServerError)
		return
	}
	respone := struct {
		Status string         `json:"status"`
		Count  int            `json:"count"`
		Data   []teacher.Exec `json:"data"`
	}{
		Status: "Success",
		Count:  len(teachersList),
		Data:   teachersList,
	}
	fmt.Println(respone)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respone)

}

func GetOneExecHandler(w http.ResponseWriter, r *http.Request) {

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

	var teacher teacher.Exec
	err = db.QueryRow("SELECT id, first_name, last_name, email, username, password, user_created_at, password_changed_at, password_reset_code, password_code_expires,inactive_status, role FROM execs WHERE id = ?", id).Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Username, &teacher.Password,
		&teacher.UserCreatedAt, &teacher.PasswordChangedAt, &teacher.PasswordResetCode, &teacher.PasswordCodeExpires, &teacher.Role)
	if err == sql.ErrNoRows {
		http.Error(w, "Exec not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database query Error ", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(teacher)

}

func AddExecsHandler(w http.ResponseWriter, r *http.Request) {
	// mutex.Lock()
	// defer mutex.Unlock()

	err, addedTeachers, ok := sqlconnect.AddExecsHandlerDB(w, r)
	if !ok {
		http.Error(w, "Error for PatchExecsHandlerDB", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := struct {
		Status string         `json:"status"`
		Count  int            `json:"count"`
		Data   []teacher.Exec `json:"data"`
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

func PatchExecsHandler(w http.ResponseWriter, r *http.Request) {

	showUpdatedTeachers, ok := sqlconnect.PatchExecsHandlerDB(w, r)
	if !ok {
		http.Error(w, "Error for PatchExecsHandlerDB", http.StatusInternalServerError)
		return
	}

	respone := struct {
		Status string         `json:"status"`
		Count  int            `json:"count"`
		Data   []teacher.Exec `json:"data"`
	}{
		Status: "Success",
		Count:  len(showUpdatedTeachers),
		Data:   showUpdatedTeachers,
	}

	json.NewEncoder(w).Encode(respone)
}

func DeleteExecHandler(w http.ResponseWriter, r *http.Request) {

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
	result, err := db.Exec("DELETE FROM execs WHERE id=?", id)
	if err != nil {

		http.Error(w, "Error deleting from db", http.StatusInternalServerError)
		return
	}

	rowsAffect, err := result.RowsAffected()
	if err != nil {
		return
	}
	if rowsAffect == 0 {
		http.Error(w, "Exec Not Found", http.StatusNotFound)
		return
	}

	response := struct {
		Status string `json:"status"`
	}{
		Status: "Successfully Deleted",
	}

	json.NewEncoder(w).Encode(response)
}

func LoginExecsHandler(w http.ResponseWriter, r *http.Request) {
	var user teacher.Exec
	json.NewDecoder(r.Body).Decode(&user)

	username := user.Username
	inputPassword := user.Password

	msg, tokenString, ok := sqlconnect.LoginExecsHandlerDB(w, username, inputPassword, user)
	if !ok {
		http.Error(w, "Error for LoginExecsHandlerDB", http.StatusInternalServerError)
		return
	}

	response := struct {
		Token string `json:"token"`
	}{
		Token: tokenString,
	}

	fmt.Println(msg)
	// response := struct {
	// 	Message string `json:"message"`
	// }{
	// 	Message: msg,
	// }
	json.NewEncoder(w).Encode(response)

}

func LogoutExecsHandler(w http.ResponseWriter, r *http.Request) {

	http.SetCookie(w, &http.Cookie{
		Name:     "Bearer",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(10 * time.Second),
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Successfully Logged out")}`))
}

func UpdatePasswordExecHandler(w http.ResponseWriter, r *http.Request) {

	// var upPass struct{
	// 	username string `json:"username"`
	// 	currentPasswor string `json:"currentpassword"`
	// 	newPassword string `json:"newpassword"`
	// }
	var upPass models.UpdatePasswordReq
	fmt.Println("UpPass:", upPass)
	json.NewDecoder(r.Body).Decode(&upPass)

	ok := sqlconnect.UpdatePasswordExecHandlerDB(w, upPass)
	if !ok {
		http.Error(w, "Error for UpdatePasswordExecHandlerDB", http.StatusInternalServerError)
		return
	}

	respone := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}{
		Status:  "Success",
		Message: "Password successfully updated",
	}

	json.NewEncoder(w).Encode(respone)

}

// func GenerateInsertQueryforExec(table string, model interface{}) string {
// 	modelType := reflect.TypeOf(model)
// 	if modelType.Kind() == reflect.Ptr {
// 		modelType = modelType.Elem()
// 	}

// 	var columns, placeholders string

// 	for i := 0; i < modelType.NumField(); i++ {
// 		dbTag := modelType.Field(i).Tag.Get("db")
// 		dbTag = strings.TrimSuffix(dbTag, ",omitempty")
// 		if dbTag != "" && dbTag != "id" {
// 			if columns != "" {
// 				columns += ", "
// 				placeholders += ", "
// 			}
// 			columns += dbTag
// 			placeholders += "?"
// 		}
// 	}

// 	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s)", table, columns, placeholders)
// 	fmt.Println("Query:", query)
// 	return query
// }
