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
)

func GetTeachersHandlerDB(w http.ResponseWriter, r *http.Request) ([]teacher.Teacher, bool) {
	db, err := ConnectDb()
	if err != nil {
		http.Error(w, "Error connecting to db", http.StatusInternalServerError)
		return nil, false
	}
	defer db.Close()

	query := "SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE 1=1"
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
		return nil, false
	}
	defer rows.Close()
	teachersList := make([]teacher.Teacher, 0)
	//teachersList := []teacher.Teacher{}

	for rows.Next() {
		var teacher teacher.Teacher
		err := rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
		if err != nil {
			http.Error(w, "Database scanning error", http.StatusInternalServerError)
			return nil, false
		}
		teachersList = append(teachersList, teacher)

	}
	return teachersList, true
}

func AddTeachersHandlerDB(w http.ResponseWriter, r *http.Request) (error, []teacher.Teacher, bool) {
	db, err := ConnectDb()
	if err != nil {
		http.Error(w, "Error connecting to db", http.StatusInternalServerError)
		return nil, nil, false
	}
	defer db.Close()

	// var rawTeachers []map[string]string{}
	rawTeachers := make([]map[string]string, 0)
	var newTeachers []teacher.Teacher
	err = json.NewDecoder(r.Body).Decode(&rawTeachers)
	if err != nil {

		fmt.Println("For decoding Error")
		return nil, nil, false
	}
	fmt.Println("Raw teachrs: ", rawTeachers)

	for _, raw := range rawTeachers {

		t := teacher.Teacher{
			FirstName: raw["first_name"],
			LastName:  raw["last_name"],
			Email:     raw["email"],
			Class:     raw["class"],
			Subject:   raw["subject"],
		}
		newTeachers = append(newTeachers, t)
	}

	fmt.Println("Raw Teacher: ", rawTeachers)
	fmt.Println("New Teacher: ", newTeachers)

	var teach teacher.Teacher
	stmt, err := db.Prepare(utils.GenerateInsertQuery("teachers", teach))
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

	typ := reflect.TypeOf(teacher.Teacher{})

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

	addedTeachers := make([]teacher.Teacher, len(newTeachers))

	for i, newTeacher := range newTeachers {

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

func PatchTeachersHandlerDB(w http.ResponseWriter, r *http.Request) ([]teacher.Teacher, bool) {
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

	var showUpdatedTeachers []teacher.Teacher

	for _, update := range updates {
		idStr, ok := update["id"].(string)
		if !ok {
			fmt.Println(ok)
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Error converting ID into int", http.StatusInternalServerError)
			return nil, false
		}
		var teacherDb teacher.Teacher

		if err != nil {
			tx.Rollback()
			http.Error(w, "Transaction err", http.StatusInternalServerError)
			return nil, false
		}
		rows := tx.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ? ", id)
		err = rows.Scan(&teacherDb.ID, &teacherDb.FirstName, &teacherDb.LastName, &teacherDb.Email, &teacherDb.Class, &teacherDb.Subject)
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

		_, err = tx.Exec("UPDATE teachers SET first_name = ? , last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?", teacherDb.FirstName, teacherDb.LastName, teacherDb.Email, teacherDb.Class, teacherDb.Subject, teacherDb.ID)
		if err != nil {
			http.Error(w, "Executing err", http.StatusInternalServerError)
			return nil, false
		}
		showUpdatedTeachers = append(showUpdatedTeachers, teacherDb)
	}

	tx.Commit()
	return showUpdatedTeachers, true
}
