package handlers

import (
	"indexer/internal/database"
	"indexer/internal/session"
	"log"
	"net/http"
)

type inputSearch struct {
	SearchTerm string `json:"searchTerm"`
}

type PageData struct {
	Username    string
	Entries     []database.Data
	CurrentPage int
	TotalPages  int
	Query       string
}

func RootHandler(db *database.DBServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data []database.Data

		page, searchTerm := getPaginationParams(r)

		offset, totalPages, err := handlePagination(db, w, page, searchTerm)
		if err != nil {
			return
		}

		data, err = db.GetData(nil, &searchTerm, offset)

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

		err = rootTemplates.ExecuteTemplate(w, "index.html", pData)

		if err != nil {
			log.Println("[root.go] template error:", err.Error())
			return
		}
	}
}
