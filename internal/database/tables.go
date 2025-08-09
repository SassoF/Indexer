package database

import "database/sql"

func createTables(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS user (
			username VARCHAR(50) NOT NULL PRIMARY KEY,
			password VARCHAR(255) NOT NULL,
			email VARCHAR(100) NOT NULL UNIQUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP()
		)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS data (
			hash VARCHAR(255) NOT NULL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			category VARCHAR(50) NOT NULL,
			description TEXT NOT NULL DEFAULT '',
			time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP(),
			size VARCHAR(10) NOT NULL,
			uploader VARCHAR(50) NOT NULL,
			INDEX (uploader),
			FOREIGN KEY (uploader) REFERENCES user(username)
		)`)
	if err != nil {
		return err
	}

	return nil
}
