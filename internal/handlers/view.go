package handlers

import (
	"database/sql"
	"html/template"
	"indexer/internal/database"
	"indexer/internal/session"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/anacrolix/torrent/metainfo"
)

func isValidSHA1(hash string) bool {
	matched, err := regexp.MatchString(`^[a-fA-F0-9]{40}$`, hash)
	return err == nil && matched
}

type View struct {
	Username    string
	Magnet      template.URL
	Data        *database.Data
	Description template.HTML
	TreeRoot    *TreeNode
}

func ViewHandler(db *database.DBServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if ok := isValidSHA1(r.PathValue("hash")); !ok {
			http.Error(w, `{"error": "The hash is not valid"}`, http.StatusBadRequest)
			return
		}

		data, err := db.GetFile(r.PathValue("hash"))

		if err != nil {
			switch err {
			case sql.ErrNoRows:
				http.Error(w, "404 file not found", http.StatusNotFound)
				return
			default:
				log.Println("[view.go]", err.Error())
				http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
				return
			}
		}

		metaInfo, err := metainfo.LoadFromFile("/app/uploads/" + data.Hash + ".torrent")
		if err != nil {
			switch err.(type) {
			case *os.PathError:
				http.Error(w, "404 file not found", http.StatusNotFound)
				return
			default:
				http.Error(w, `{"error": "Invalid torrent file"}`, http.StatusBadRequest)
				return
			}
		}

		view := &View{Data: data}
		if view.Data.Description != "" {
			view.Description = template.HTML(view.Data.Description)
		}

		info, err := metaInfo.UnmarshalInfo()
		if err != nil {
			http.Error(w, `{"error": "Invalid torrent file"}`, http.StatusBadRequest)
			return
		}

		if magnet, err := metaInfo.MagnetV2(); err != nil {
			log.Println("[upload.go]", err.Error())
			http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
			return
		} else {
			view.Magnet = template.URL(magnet.String())
		}

		view.Username = ""

		if u := session.SessionManager.Get(r.Context(), "username"); u != nil {
			if u, ok := u.(string); ok {
				view.Username = u
			}
		}

		view.TreeRoot = buildTree(&info)

		err = viewTemplate.ExecuteTemplate(w, "view.html", *view)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	}
}
