package handlers

import (
	"indexer/internal/database"
	"indexer/internal/session"
	"net/http"
)

func DefineRoutes(mux *http.ServeMux, dtbs *database.DBServer) {

	mux.Handle("GET /static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("/app/web/static"))),
	)
	mux.Handle("GET /download/",
		http.StripPrefix("/download/", http.FileServer(http.Dir("/app/uploads/"))),
	)

	mux.Handle("POST /login", LoginHandler(dtbs))
	mux.Handle("POST /register", RegisterHandler(dtbs))
	mux.Handle("GET /logout", http.HandlerFunc(LogoutHandler))

	mux.HandleFunc("GET /profile/", session.RedirectIfLogged(UserHandler))
	mux.HandleFunc("GET /profile/{profile}", ProfileHandler(dtbs))

	mux.HandleFunc("GET /upload", session.RequireAuth(UploadPageHandler))
	mux.Handle("POST /upload", session.RequireAuth(UploadFileHandler(dtbs)))

	mux.HandleFunc("GET /view/{hash}", ViewHandler(dtbs))

	mux.Handle("GET /{$}", RootHandler(dtbs))
}
