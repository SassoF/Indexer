package handlers

import (
	"indexer/internal/database"
	"indexer/internal/session"
	"log"
	"net/http"
)

func ProfileHandler(db *database.DBServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data []database.Data

		page, searchTerm := getPaginationParams(r)

		offset, totalPages, err := handlePagination(db, w, page, searchTerm)
		if err != nil {
			return
		}

		profile := r.PathValue("profile")
		data, err = db.GetData(&profile, &searchTerm, offset)

		if err != nil {
			log.Println("[root.go] Database error:", err.Error())
			http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
			return
		}

		username := ""
		if u := session.SessionManager.Get(r.Context(), "username"); u != nil {
			if u, ok := u.(string); ok {
				username = u
			}
		}

		pData := PageData{
			Username:    username,
			Entries:     data,
			CurrentPage: page,
			TotalPages:  totalPages,
			Query:       searchTerm,
		}

		err = rootTemplates.ExecuteTemplate(w, "profile.html", pData)

		if err != nil {
			log.Println("[root.go] template error:", err.Error())
			return
		}
	}

}
