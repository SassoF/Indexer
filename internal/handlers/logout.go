package handlers

import (
	"indexer/internal/session"
	"log"
	"net/http"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	err := session.SessionManager.Destroy(r.Context())
	if err != nil {
		log.Println("[logout.go]", err.Error())
		http.Error(w, `{"error": "Internal status error"}`, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
