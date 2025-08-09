package database

func (s *DBServer) CountPage(page *int, name *string) (int, int, error) {

	var totalRecords int
	var err error

	if *name != "" {
		err = s.DB.QueryRow(
			"SELECT COUNT(*) FROM data WHERE name LIKE ?",
			"%"+*name+"%").Scan(&totalRecords)
	} else {
		err = s.DB.QueryRow("SELECT COUNT(*) FROM data").Scan(&totalRecords)
	}

	if err != nil {
		return 0, 0, err
	}

	totalPages := totalRecords / 50
	if totalRecords%50 != 0 {
		totalPages++
	}

	if *page > totalPages {
		*page = totalPages
	}

	return (*page - 1) * 50, totalPages, nil

}
