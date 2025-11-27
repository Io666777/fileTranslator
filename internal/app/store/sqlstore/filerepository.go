package sqlstore

import (
	"database/sql"

	"github.com/Io666777/fileTranslator/internal/app/model"
	"github.com/Io666777/fileTranslator/internal/app/store"
)

type FileRepository struct {
	store *Store
}

func (r *FileRepository) Create(f *model.File) error {
	if err := f.Validate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO files (user_id, filename, original_path, file_size, mime_type, status) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		f.UserID,
		f.Filename,
		f.OriginalPath,
		f.FileSize,
		f.MimeType,
		f.Status,
	).Scan(&f.ID)
}

func (r *FileRepository) Find(id int) (*model.File, error) {
	f := &model.File{}
	err := r.store.db.QueryRow(
		"SELECT id, user_id, filename, original_path, file_size, mime_type, status, created_at FROM files WHERE id = $1",
		id,
	).Scan(&f.ID, &f.UserID, &f.Filename, &f.OriginalPath, &f.FileSize, &f.MimeType, &f.Status, &f.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	return f, nil
}

func (r *FileRepository) FindByUserID(userID int) ([]*model.File, error) {
	rows, err := r.store.db.Query(
		"SELECT id, user_id, filename, original_path, file_size, mime_type, status, created_at FROM files WHERE user_id = $1 ORDER BY created_at DESC",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []*model.File
	for rows.Next() {
		f := &model.File{}
		err := rows.Scan(&f.ID, &f.UserID, &f.Filename, &f.OriginalPath, &f.FileSize, &f.MimeType, &f.Status, &f.CreatedAt)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}

	return files, nil
}

func (r *FileRepository) Delete(id int) error {
	_, err := r.store.db.Exec("DELETE FROM files WHERE id = $1", id)
	return err
}

func (r *FileRepository) Update(f *model.File) error {
	_, err := r.store.db.Exec(
		"UPDATE files SET status = $1 WHERE id = $2",
		f.Status, f.ID,
	)
	return err
}