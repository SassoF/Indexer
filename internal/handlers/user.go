package handlers

import (
	"indexer/internal/database"
	"net/http"
	"net/mail"
	"regexp"
	"strings"
	"unicode"
)

func UserHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "/app/web/templates/auth.html")
}

func isValidEmail(email string) bool {
	if len(email) > 100 {
		return false
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}

	emailRegex := regexp.MustCompile(`^[\w\-\.]+@([\w-]+\.)+[\w-]{2,}$`)
	if !emailRegex.MatchString(email) {
		return false
	}

	return true
}

func isValidUsername(username string) bool {
	if len(username) < 3 || len(username) > 30 {
		return false
	}

	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !usernameRegex.MatchString(username) {
		return false
	}

	if strings.HasPrefix(username, "-") ||
		strings.HasPrefix(username, "_") ||
		strings.HasSuffix(username, "-") ||
		strings.HasSuffix(username, "_") {
		return false
	}

	return true
}

func isStrongPassword(password string) bool {
	if len(password) < 8 || len(password) > 30 {
		return false
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

func isValidUserRegistration(u *database.User) (string, bool) {
	if !isValidEmail(u.Email) {
		return `{"error": "Please provide a valid email address"}`, false
	}

	if !isValidUsername(u.Username) {
		return `{"error": "Username must be 3-30 characters long and can only contain letters, numbers, underscores and hyphens"}`, false
	}

	if !isStrongPassword(u.Password) {
		return `{"error": "Password must be 8-30 characters long and contain uppercase, lowercase, number and special character"}`, false
	}

	return "", true
}

func isValidUserLogin(u *database.User) (string, bool) {
	if !isValidUsername(u.Username) {
		return `{"error": "Invalid username"}`, false
	}

	if u.Password == "" || len(u.Password) > 255 {
		return `{"error": "Invalid password"}`, false
	}

	return "", true
}
