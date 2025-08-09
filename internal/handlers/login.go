package handlers

import (
	"database/sql"
	"indexer/internal/database"
	"indexer/internal/session"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(db *database.DBServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		u, err := getPostFormValue(r)
		if err != nil {
			http.Error(w, `{"error": "Error parsing form"}`, http.StatusBadRequest)
			return
		}

		if msg, ok := isValidUserLogin(u); !ok {
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		err = db.LoginUser(u)

		if err != nil {
			switch err {
			case sql.ErrNoRows:
				http.Error(w, `{"error": "The user does not exist."}`, http.StatusBadRequest)
				return
			case bcrypt.ErrMismatchedHashAndPassword:
				http.Error(w, `{"error": "The password is incorrect."}`, http.StatusBadRequest)
				return
			default:
				log.Println("[login.go]", err.Error())
				http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
				return
			}
		}

		err = session.SessionManager.RenewToken(r.Context())
		if err != nil {
			log.Println("[login.go]", err.Error())
			http.Error(w, `{"error": "Internal status error"}`, http.StatusInternalServerError)
			return
		}

		session.SessionManager.Put(r.Context(), "username", u.Username)

		http.Redirect(w, r, "/", http.StatusSeeOther)

	}
}
