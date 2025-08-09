package database

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string
	Password string
	Email    string
}

func (s *DBServer) RegisterUser(u *User) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	result, err := s.DB.Exec(
		"INSERT INTO user (username, password, email) VALUES (?, ?, ?)",
		u.Username,
		string(hashedPassword),
		u.Email,
	)

	if err != nil {
		return err
	}

	_, err = result.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}
