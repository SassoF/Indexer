package handlers

import (
	"indexer/internal/database"
	"indexer/internal/session"
	"log"
	"net/http"
	"strings"

	"github.com/go-sql-driver/mysql"
)

func getPostFormValue(r *http.Request) (*database.User, error) {

	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	u := new(database.User)
	u.Username = r.PostFormValue("username")
	u.Password = r.PostFormValue("password")
	return u, nil
}

func RegisterHandler(db *database.DBServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		u, err := getPostFormValue(r)
		if err != nil {
			http.Error(w, `{"error": "Error parsing form"}`, http.StatusBadRequest)
			return
		}
		u.Email = r.PostFormValue("email")

		if msg, ok := isValidUserRegistration(u); !ok {
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		err = db.RegisterUser(u)
		if err != nil {
			if mysqlErr, ok := err.(*mysql.MySQLError); ok {
				switch mysqlErr.Number {
				case 1062:
					msg := `{"error": "The value is already registered."}`
					if strings.Contains(mysqlErr.Message, "for key ") {
						parts := strings.Split(mysqlErr.Message, "for key ")
						if len(parts) > 1 {
							switch strings.Trim(parts[1], "'") {
							case "PRIMARY":
								msg = `{"error": "The username already exists."}`
							case "email":
								msg = `{"error": "The email is already registered."}`
							}
						}
					}

					http.Error(w, msg, http.StatusBadRequest)
					return
				}
			}
			log.Println("[register.go]", err.Error())
			http.Error(w, `{"error": "Internal status error"}`, http.StatusInternalServerError)
			return
		}

		err = session.SessionManager.RenewToken(r.Context())
		if err != nil {
			log.Println("[register.go]", err.Error())
			http.Error(w, `{"error": "Internal status error"}`, http.StatusInternalServerError)
			return
		}

		session.SessionManager.Put(r.Context(), "username", u.Username)

		http.Redirect(w, r, "/", http.StatusSeeOther)

	}
}
