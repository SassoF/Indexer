package database

import (
	"database/sql"
	"strings"
)

type DBServer struct {
	DB *sql.DB
}

func (s *DBServer) GetData(username *string, searchTerm *string, page int) ([]Data, error) {

	var queryParts []string
	var args []any

	if username != nil {
		queryParts = append(queryParts, "uploader = ?")
		args = append(args, *username)
	}

	if searchTerm != nil && *searchTerm != "" {
		queryParts = append(queryParts, "name LIKE ?")
		args = append(args, "%"+*searchTerm+"%")
	}

	whereClause := ""
	if len(queryParts) > 0 {
		whereClause = "WHERE " + strings.Join(queryParts, " AND ")
	}
	args = append(args, page)

	rows, err := s.DB.Query(
		"SELECT hash, name, category, time, size, uploader FROM data "+
			whereClause+
			" ORDER BY time DESC LIMIT 50 OFFSET ?",
		args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	results := make([]Data, 0, 50)

	for rows.Next() {
		r := new(Data)
		err := rows.Scan(&r.Hash, &r.Name, &r.Category, &r.Time, &r.Size, &r.Uploader)

		if err != nil {
			return nil, err
		}
		results = append(results, *r)
	}
	return results, nil
}
