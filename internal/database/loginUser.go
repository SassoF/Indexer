package database

import (
	"golang.org/x/crypto/bcrypt"
)

func (s *DBServer) LoginUser(u *User) error {

	var hashedPassword string

	err := s.DB.QueryRow(
		"SELECT password FROM user WHERE username = ?",
		u.Username,
	).Scan(&hashedPassword)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(u.Password))
	if err != nil {
		return err
	}

	return nil
}
