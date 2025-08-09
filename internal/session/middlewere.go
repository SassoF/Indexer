package session

import (
	"net/http"
)

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if username := SessionManager.Get(r.Context(), "username"); username == nil {
			http.Redirect(w, r, "/profile", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func RedirectIfLogged(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if username := SessionManager.Get(r.Context(), "username"); username != nil {
			if usernameStr, ok := username.(string); ok {
				http.Redirect(w, r, "/profile/"+usernameStr, http.StatusSeeOther)
				return
			}
		}
		next.ServeHTTP(w, r)
	}
}
