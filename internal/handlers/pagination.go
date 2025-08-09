package handlers

import (
	"indexer/internal/database"
	"log"
	"net/http"
	"strconv"
)

func getPaginationParams(r *http.Request) (int, string) {
	searchTerm := r.URL.Query().Get("q")

	page, err := strconv.Atoi(r.URL.Query().Get("p"))
	if err != nil || page < 1 {
		page = 1
	}

	return page, searchTerm
}

func handlePagination(db *database.DBServer, w http.ResponseWriter, page int, searchTerm string) (int, int, error) {
	offset, totalPages, err := db.CountPage(&page, &searchTerm)
	if err != nil {
		log.Println("[profile.go] Database error:", err.Error())
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return 0, 0, err
	}
	return offset, totalPages, nil
}
