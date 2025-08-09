package database

func (s *DBServer) UploadFile(f *Data) error {
	result, err := s.DB.Exec(
		"INSERT INTO data (hash, name, category, description, size, uploader) VALUES (?, ?, ?, ?, ?, ?)",
		f.Hash,
		f.Name,
		f.Category,
		f.Description,
		f.Size,
		f.Uploader,
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
