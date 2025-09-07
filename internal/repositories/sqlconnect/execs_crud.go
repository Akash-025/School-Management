package sqlconnect

import (
	"database/sql"
	"net/http"
	teacher "restapi/internal/models"
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
