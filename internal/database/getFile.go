package database

func (s *DBServer) GetFile(hash string) (*Data, error) {

	data := new(Data)

	err := s.DB.QueryRow(
		"SELECT hash, name, category, description, time, size, uploader FROM data WHERE hash = ?",
		hash,
	).Scan(&data.Hash,
		&data.Name,
		&data.Category,
		&data.Description,
		&data.Time,
		&data.Size,
		&data.Uploader,
	)
	if err != nil {
		return nil, err
	}

	return data, nil
}
