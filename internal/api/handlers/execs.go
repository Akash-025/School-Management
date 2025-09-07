package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"restapi/internal/models"
	teacher "restapi/internal/models"
	"restapi/internal/repositories/sqlconnect"
	"restapi/pkg/utils"
	"strconv"
	"strings"
	"time"
)

func GetExecsHandler(w http.ResponseWriter, r *http.Request) {

	db, err := sqlconnect.ConnectDb()
	if err != nil {
		http.Error(w, "Error connecting to db", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	query := "SELECT id, first_name, last_name, email, username, password, user_created_at, password_changed_at, password_reset_code, password_code_expires, inactive_status, role FROM execs WHERE 1=1"
	var args []interface{}

	params := map[string]string{
		"first_name":            "first_name",
		"last_name":             "last_name",
		"email":                 "email",
		"username":              "username",
		"password":              "password",
		"user_created_at":       "user_created_at",
		"password_changed_at":   "password_changed_at",
		"password_reset_code":   "password_reset_code",
		"password_code_expires": "password_code_expires",
		"inactive_status":       "inactive_status",
		"role":                  "role",
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

	fmt.Println("query:", query, "Arg", args)
	rows, err := db.Query(query, args...)
	if err != nil {
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}
	fmt.Println("Rows1:", rows)
	defer rows.Close()
	teachersList := make([]teacher.Exec, 0)
	//teachersList := []teacher.Exec{}

	for rows.Next() {
		var teacher teacher.Exec
		err := rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Username, &teacher.Password,
			&teacher.UserCreatedAt, &teacher.PasswordChangedAt, &teacher.PasswordResetCode, &teacher.PasswordCodeExpires, &teacher.InactiveStatus, &teacher.Role)
		fmt.Println("Rows:", teacher)
		if err != nil {
			http.Error(w, "Database scanning error", http.StatusInternalServerError)
			return
		}
		teachersList = append(teachersList, teacher)

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

	db, err := sqlconnect.ConnectDb()
	if err != nil {
		http.Error(w, "Error connecting to db", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// var rawTeachers []map[string]string{}
	rawTeachers := make([]map[string]string, 0)
	var newTeachers []teacher.Exec
	err = json.NewDecoder(r.Body).Decode(&rawTeachers)
	if err != nil {

		fmt.Println("For decoding Error")
		return
	}
	fmt.Println("Raw execs: ", rawTeachers)

	for _, raw := range rawTeachers {

		t := teacher.Exec{
			FirstName:           raw["first_name"],
			LastName:            raw["last_name"],
			Email:               raw["email"],
			Username:            raw["username"],
			Password:            raw["password"],
			UserCreatedAt:       sql.NullString{String: raw["user_created_at"], Valid: raw["user_created_at"] != ""},
			PasswordChangedAt:   sql.NullString{String: raw["password_changed_at"], Valid: raw["password_changed_at"] != ""},
			PasswordResetCode:   sql.NullString{String: raw["password_reset_code"], Valid: raw["password_reset_code"] != ""},
			PasswordCodeExpires: sql.NullString{String: raw["password_code_expires"], Valid: raw["password_code_expires"] != ""},
			InactiveStatus:      raw["inactive_status"] == "true",
			Role:                raw["role"],
		}

		newTeachers = append(newTeachers, t)
	}

	fmt.Println("Raw Teacher: ", rawTeachers)
	fmt.Println("New Teacher: ", newTeachers)

	var teach teacher.Exec
	stmt, err := db.Prepare(GenerateInsertQueryforExec("execs", &teach))
	if err != nil {
		http.Error(w, "Error to prepare db connection", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// All Fields are required
	for _, teacher := range newTeachers {

		val := reflect.ValueOf(teacher)
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			if field.String() == "" {
				http.Error(w, "All fields are required", http.StatusBadRequest)
				return
			}
		}
	}

	// Extra Fields are not allowed

	typ := reflect.TypeOf(teacher.Exec{})

	fields := []string{}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i).Tag.Get("json")
		f := strings.TrimSuffix(field, ",omitempty")
		fields = append(fields, f)
	}
	allowedFields := make(map[string]struct{})
	for _, field := range fields {
		allowedFields[field] = struct{}{}
	}

	fmt.Println("Fields are:", fields)
	for _, teach := range rawTeachers {
		for v := range teach {
			_, ok := allowedFields[v]
			fmt.Println("Allowed fields: ", v)
			if !ok {
				http.Error(w, "Extra fields are not allowed", http.StatusBadRequest)
				return
			}
		}
	}

	addedTeachers := make([]teacher.Exec, len(newTeachers))

	for i, newTeacher := range newTeachers {

		// Password Hashing
		if newTeacher.Password != "" {
			pass, err := utils.HashPassword(newTeacher.Password, w)
			if err != nil {
				return
			}
			newTeacher.Password = pass
		}

		values := utils.GetStructValues(newTeacher)
		res, err := stmt.Exec(values...)
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

	db, err := sqlconnect.ConnectDb()
	if err != nil {
		http.Error(w, "Error connecting to db", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var updates []map[string]interface{}

	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Decoding err", http.StatusInternalServerError)
		return
	}
	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Transaction starting err", http.StatusInternalServerError)
		return
	}

	var showUpdatedTeachers []teacher.Exec

	for _, update := range updates {
		idStr, ok := update["id"].(string)
		if !ok {
			fmt.Println("Error for id:", ok)
		}

		id, err := strconv.Atoi(idStr)
		//id := int64(idr)
		if err != nil {
			http.Error(w, "Error converting ID into int", http.StatusInternalServerError)
			return
		}
		var teacherDb teacher.Exec
		fmt.Println("Valo ID:", id)

		if err != nil {
			tx.Rollback()
			http.Error(w, "Transaction err", http.StatusInternalServerError)
			return
		}
		rows := tx.QueryRow("SELECT id, first_name, last_name, email, username, password, user_created_at, password_changed_at, password_reset_code, password_code_expires, inactive_status, role FROM execs WHERE id = ? ", id)
		err = rows.Scan(&teacherDb.ID, &teacherDb.FirstName, &teacherDb.LastName, &teacherDb.Email, &teacherDb.Username, &teacherDb.Password,
			&teacherDb.UserCreatedAt, &teacherDb.PasswordChangedAt, &teacherDb.PasswordResetCode, &teacherDb.PasswordCodeExpires, &teacherDb.InactiveStatus, &teacherDb.Role)
		if err != nil {
			tx.Rollback()
			if err == sql.ErrNoRows {
				http.Error(w, "Teacher Not Found", http.StatusNotFound)
				return
			}
			http.Error(w, "Scanning err", http.StatusInternalServerError)
			return
		}

		teacherVal := reflect.ValueOf(&teacherDb).Elem()
		teacherType := teacherVal.Type()

		for k, v := range update {
			if k == "id" {
				continue
			}
			for i := 0; i < teacherVal.NumField(); i++ {

				field := teacherType.Field(i)
				kk := strings.TrimSuffix(field.Tag.Get("json"), ",omitempty")
				if k == kk {
					fieldVal := teacherVal.Field(i)
					if fieldVal.CanSet() {
						val := reflect.ValueOf(v)
						if val.Type().ConvertibleTo(fieldVal.Type()) {
							fieldVal.Set(val.Convert(fieldVal.Type()))
						} else {
							tx.Rollback()
							fmt.Println("Connot convert")
							return
						}
					}
					break

				}
			}
		}

		_, err = tx.Exec("UPDATE execs SET first_name = ?, last_name = ?, email = ?, username = ?, password = ?, user_created_at = ?, password_changed_at = ?, password_reset_code = ?, password_code_expires = ?, inactive_status = ?, role = ? WHERE id = ?",
			teacherDb.FirstName, teacherDb.LastName, teacherDb.Email, teacherDb.Username, teacherDb.Password, teacherDb.UserCreatedAt,
			teacherDb.PasswordChangedAt, teacherDb.PasswordResetCode, teacherDb.PasswordCodeExpires, teacherDb.InactiveStatus, teacherDb.Role, teacherDb.ID)

		if err != nil {
			http.Error(w, "Executing err", http.StatusInternalServerError)
			return
		}
		showUpdatedTeachers = append(showUpdatedTeachers, teacherDb)
	}

	tx.Commit()

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

	userFromDb, NotOk := sqlconnect.LoginDb(w, username)
	if NotOk {
		return
	}

	// Password checking
	msg, NotOk := utils.PasswordCheck(userFromDb.Password, w, inputPassword)
	if NotOk {
		http.Error(w, "Password is incorrect", http.StatusForbidden)
		return
	}
	// Cookie
	tokenString, err := utils.SignToken(user.ID, user.Username, user.Role)
	if err != nil {
		return 
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "Bearer",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(100 * time.Second),
		SameSite: http.SameSiteStrictMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "Test",
		Value:    "Testing",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(10 * time.Second),
		SameSite: http.SameSiteStrictMode,
	})

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

func LogoutExecsHandler(w http.ResponseWriter, r *http.Request){

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

func UpdatePasswordExecHandler(w http.ResponseWriter, r *http.Request){

	// var upPass struct{
	// 	username string `json:"username"`
	// 	currentPasswor string `json:"currentpassword"`
	// 	newPassword string `json:"newpassword"`
	// }
	var upPass models.UpdatePasswordReq
	fmt.Println("UpPass:", upPass)
	json.NewDecoder(r.Body).Decode(&upPass)

	db, err := sqlconnect.ConnectDb()
	if err != nil {
		http.Error(w, "Error connecting to db", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var execs teacher.Exec
	db.QueryRow("SELECT password FROM execs WHERE username = ?", upPass.Username).Scan(&execs.Password)

	passwordFromDb := execs.Password

	_, NotOk := utils.PasswordCheck(passwordFromDb, w, upPass.CurrentPassword)
	if NotOk {
		http.Error(w, "The password you entered is incorrect", http.StatusForbidden)
		return
	}

	hashedPassword, err := utils.HashPassword(upPass.NewPassword, w)
	if err != nil {
		http.Error(w, "Password hashing error", http.StatusInternalServerError)
		return 
	}

	_, err = db.Exec("UPDATE execs SET password = ? WHERE username = ?", hashedPassword, upPass.Username)
	if err != nil {
		http.Error(w, "Password update query error", http.StatusInternalServerError)
		return
	}

	respone := struct{
		Status string `json:"status"`
		Message string `json:"message"`
	}{
		Status: "Success",
		Message: "Password successfully updated",
	}

	json.NewEncoder(w).Encode(respone)

}

func GenerateInsertQueryforExec(table string, model interface{}) string {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	var columns, placeholders string

	for i := 0; i < modelType.NumField(); i++ {
		dbTag := modelType.Field(i).Tag.Get("db")
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

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s)", table, columns, placeholders)
	fmt.Println("Query:", query)
	return query
}
