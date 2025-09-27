package sqlconnect

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	teacher "restapi/internal/models"
	"restapi/pkg/utils"
	"strconv"
	"strings"
	"time"
)

func LoginDb(w http.ResponseWriter, username string) (teacher.Exec, bool) {
	db, err := ConnectDb()
	if err != nil {
		http.Error(w, "Error connecting to db", http.StatusInternalServerError)
		return teacher.Exec{}, true
	}
	defer db.Close()
	var userFromDb teacher.Exec
	err = db.QueryRow("SELECT password FROM execs WHERE username = ?", username).Scan(&userFromDb.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Execs not exist", http.StatusInternalServerError)
			return teacher.Exec{}, true
		}
		http.Error(w, "DB query error", http.StatusInternalServerError)
		return teacher.Exec{}, true
	}
	return userFromDb, false
}

func GetExecsHandlerDB(w http.ResponseWriter, r *http.Request) ([]teacher.Exec, bool) {
	db, err := ConnectDb()
	if err != nil {
		http.Error(w, "Error connecting to db", http.StatusInternalServerError)
		return nil, false
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
		return nil, false
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
			return nil, false
		}
		teachersList = append(teachersList, teacher)

	}
	return teachersList, true
}

func AddExecsHandlerDB(w http.ResponseWriter, r *http.Request) (error, []teacher.Exec, bool) {
	db, err := ConnectDb()
	if err != nil {
		http.Error(w, "Error connecting to db", http.StatusInternalServerError)
		return nil, nil, false
	}
	defer db.Close()

	// var rawTeachers []map[string]string{}
	rawTeachers := make([]map[string]string, 0)
	var newTeachers []teacher.Exec
	err = json.NewDecoder(r.Body).Decode(&rawTeachers)
	if err != nil {

		fmt.Println("For decoding Error")
		return nil, nil, false
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
	stmt, err := db.Prepare(utils.GenerateInsertQuery("execs", &teach))
	if err != nil {
		http.Error(w, "Error to prepare db connection", http.StatusInternalServerError)
		return nil, nil, false
	}
	defer stmt.Close()

	// All Fields are required
	for _, teacher := range newTeachers {

		val := reflect.ValueOf(teacher)
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			if field.String() == "" {
				http.Error(w, "All fields are required", http.StatusBadRequest)
				return nil, nil, false
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
				return nil, nil, false
			}
		}
	}

	addedTeachers := make([]teacher.Exec, len(newTeachers))

	for i, newTeacher := range newTeachers {

		// Password Hashing
		if newTeacher.Password != "" {
			pass, err := utils.HashPassword(newTeacher.Password, w)
			if err != nil {
				return nil, nil, false
			}
			newTeacher.Password = pass
		}

		values := utils.GetStructValues(newTeacher)
		res, err := stmt.Exec(values...)
		if err != nil {
			http.Error(w, "Error to execute db operation", http.StatusInternalServerError)
			return nil, nil, false
		}

		lastID, err := res.LastInsertId()
		if err != nil {
			http.Error(w, "Error for lastID", http.StatusInternalServerError)
			return nil, nil, false
		}
		newTeacher.ID = int(lastID)
		addedTeachers[i] = newTeacher

	}
	return err, addedTeachers, true
}

func PatchExecsHandlerDB(w http.ResponseWriter, r *http.Request) ([]teacher.Exec, bool) {
	db, err := ConnectDb()
	if err != nil {
		http.Error(w, "Error connecting to db", http.StatusInternalServerError)
		return nil, false
	}
	defer db.Close()

	var updates []map[string]interface{}

	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Decoding err", http.StatusInternalServerError)
		return nil, false
	}
	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Transaction starting err", http.StatusInternalServerError)
		return nil, false
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
			return nil, false
		}
		var teacherDb teacher.Exec
		fmt.Println("Valo ID:", id)

		if err != nil {
			tx.Rollback()
			http.Error(w, "Transaction err", http.StatusInternalServerError)
			return nil, false
		}
		rows := tx.QueryRow("SELECT id, first_name, last_name, email, username, password, user_created_at, password_changed_at, password_reset_code, password_code_expires, inactive_status, role FROM execs WHERE id = ? ", id)
		err = rows.Scan(&teacherDb.ID, &teacherDb.FirstName, &teacherDb.LastName, &teacherDb.Email, &teacherDb.Username, &teacherDb.Password,
			&teacherDb.UserCreatedAt, &teacherDb.PasswordChangedAt, &teacherDb.PasswordResetCode, &teacherDb.PasswordCodeExpires, &teacherDb.InactiveStatus, &teacherDb.Role)
		if err != nil {
			tx.Rollback()
			if err == sql.ErrNoRows {
				http.Error(w, "Teacher Not Found", http.StatusNotFound)
				return nil, false
			}
			http.Error(w, "Scanning err", http.StatusInternalServerError)
			return nil, false
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
							return nil, false
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
			return nil, false
		}
		showUpdatedTeachers = append(showUpdatedTeachers, teacherDb)
	}

	tx.Commit()
	return showUpdatedTeachers, true
}

func LoginExecsHandlerDB(w http.ResponseWriter, username string, inputPassword string, user teacher.Exec) (string, string, bool) {
	db, err := ConnectDb()
	if err != nil {
		http.Error(w, "Error connecting to db", http.StatusInternalServerError)
		return "", "", false
	}
	defer db.Close()

	var role string
	db.QueryRow("SELECT role FROM execs WHERE username = ?", username).Scan(&role)

	userFromDb, NotOk := LoginDb(w, username)
	if NotOk {
		return "", "", false
	}

	// Password checking
	msg, NotOk := utils.PasswordCheck(userFromDb.Password, w, inputPassword)
	if NotOk {
		http.Error(w, "Password is incorrect", http.StatusForbidden)
		return "", "", false
	}
	// Cookie
	tokenString, err := utils.SignToken(user.ID, user.Username, role)
	if err != nil {
		return "", "", false
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
	return msg, tokenString, true
}

func UpdatePasswordExecHandlerDB(w http.ResponseWriter, upPass teacher.UpdatePasswordReq) bool {
	db, err := ConnectDb()
	if err != nil {
		http.Error(w, "Error connecting to db", http.StatusInternalServerError)
		return false
	}
	defer db.Close()

	var execs teacher.Exec
	db.QueryRow("SELECT password FROM execs WHERE username = ?", upPass.Username).Scan(&execs.Password)

	passwordFromDb := execs.Password

	_, NotOk := utils.PasswordCheck(passwordFromDb, w, upPass.CurrentPassword)
	if NotOk {
		http.Error(w, "The password you entered is incorrect", http.StatusForbidden)
		return false
	}

	hashedPassword, err := utils.HashPassword(upPass.NewPassword, w)
	if err != nil {
		http.Error(w, "Password hashing error", http.StatusInternalServerError)
		return false
	}

	_, err = db.Exec("UPDATE execs SET password = ? WHERE username = ?", hashedPassword, upPass.Username)
	if err != nil {
		http.Error(w, "Password update query error", http.StatusInternalServerError)
		return false
	}
	return true
}
