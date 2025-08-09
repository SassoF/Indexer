package handlers

import (
	"bytes"
	"fmt"
	"html"
	"html/template"
	"indexer/internal/database"
	"indexer/internal/session"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/anacrolix/torrent/metainfo"
	"github.com/go-sql-driver/mysql"
	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

func UploadPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("/app/web/templates/upload.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, struct {
		Username string
	}{
		Username: session.SessionManager.Get(r.Context(), "username").(string),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func hasStrictHTMLTags(input string) bool {
	p := bluemonday.StrictPolicy()
	sanitized := p.Sanitize(input)
	return html.UnescapeString(sanitized) != strings.TrimSpace(input)
}

func isTorrentFilesValid(info *metainfo.Info) bool {
	for _, file := range info.Files {
		for _, part := range file.Path {
			if hasStrictHTMLTags(part) {
				return false
			}
		}
	}
	return true
}

var categories = map[string]string{
	"1": "Movies",
	"2": "Television",
	"3": "Games",
	"4": "Music",
	"5": "Applications",
	"6": "Anime",
	"7": "Documentaries",
	"8": "Other",
	"9": "XXX",
}

func isValidCategory(category *string) bool {
	var exists bool
	*category, exists = categories[*category]
	if !exists {
		return false
	}
	return true
}

func sanitizeMarkdown(md string) (string, error) {
	markdown := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
		),
	)

	var buf bytes.Buffer
	if err := markdown.Convert([]byte(md), &buf); err != nil {
		return "", err
	}

	p := bluemonday.UGCPolicy()
	return strings.TrimSpace(p.Sanitize(buf.String())), nil
}

func UploadFileHandler(db *database.DBServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			http.Error(w, `{"error": "The file exceeds the maximum size 1MB"}`, http.StatusBadRequest)
			return
		}

		f := new(database.Data)

		f.Category = r.PostFormValue("category")
		if f.Category != "" {
			if !isValidCategory(&f.Category) {
				http.Error(w, `{"error": "Category is not valid"}`, http.StatusBadRequest)
				return
			}
		}

		f.Description = r.PostFormValue("description")
		if len(f.Description) > 2500 {
			http.Error(w, `{"error": "The description is too long"}`, http.StatusBadRequest)
			return

		}

		var err error
		f.Description, err = sanitizeMarkdown(f.Description)
		if err != nil {
			http.Error(w, `{"error": "The description is not valid"}`, http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, `{"error": "Failed to retrieve the file"}`, http.StatusBadRequest)
			return
		}
		defer file.Close()

		if header.Header.Get("Content-Type") != "application/x-bittorrent" {
			http.Error(w, `{"error": "Invalid header"}`, http.StatusBadRequest)
		}

		var buf bytes.Buffer
		tee := io.TeeReader(file, &buf)

		metaInfo, err := metainfo.Load(tee)
		if err != nil {
			http.Error(w, `{"error": "Invalid torrent file"}`, http.StatusBadRequest)
			return
		}
		info, err := metaInfo.UnmarshalInfo()
		if err != nil {
			http.Error(w, `{"error": "Invalid torrent file"}`, http.StatusBadRequest)
			return
		}

		if !isTorrentFilesValid(&info) {
			http.Error(w, `{"error": "Invalid torrent file"}`, http.StatusBadRequest)
			return
		}

		f.Name = r.PostFormValue("name")

		if f.Name == "" {
			f.Name = info.Name
		}
		if len(f.Name) > 50 {
			http.Error(w, `{"error": "Invalid name length"}`, http.StatusBadRequest)
			return
		}
		if hasStrictHTMLTags(f.Name) {
			http.Error(w, `{"error": "Invalid name"}`, http.StatusBadRequest)
			return
		}

		f.Size = formatSize(info.TotalLength())
		f.Uploader = session.SessionManager.Get(r.Context(), "username").(string)
		f.Hash = strings.ToUpper(metaInfo.HashInfoBytes().HexString())
		filename := fmt.Sprintf("/app/uploads/%s.torrent", f.Hash)

		if err := db.UploadFile(f); err != nil {
			if mysqlErr, ok := err.(*mysql.MySQLError); ok {
				switch mysqlErr.Number {
				case 1062:
					msg := `{"error": "This torrent file was already uploaded"}`
					http.Error(w, msg, http.StatusBadRequest)
					return
				}
			}
			log.Println("[upload.go]", err.Error())
			http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
			return
		}

		dst, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
		if err != nil {
			if os.IsExist(err) {
				http.Error(w, `{"error": "This torrent file was already uploaded"}`, http.StatusConflict)
			} else {
				log.Println("[upload.go]", err.Error())
				http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
			}
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, &buf); err != nil {
			http.Error(w, `{"error": "Failed to save the file"}`, http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/view/"+f.Hash, http.StatusSeeOther)
	}
}
