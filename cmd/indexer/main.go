package main

import (
	"indexer/internal/database"
	"indexer/internal/handlers"
	"indexer/internal/session"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	ip := os.Getenv("INDEXER_IP")
	port := os.Getenv("INDEXER_PORT")

	if port == "" {
		port = "8080"
	}

	session.Init()

	db, err := database.InitDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dtbs := &database.DBServer{DB: db}

	mux := http.NewServeMux()

	handlers.DefineRoutes(mux, dtbs)
	handlers.InitTemplates()

	log.Println("[main.go] Server listening on http://" + ip + ":" + port)
	log.Fatal(http.ListenAndServe(ip+":"+port, session.SessionManager.LoadAndSave(mux)))
}
