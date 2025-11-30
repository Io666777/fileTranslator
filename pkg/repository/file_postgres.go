package repository

import (
	"filetranslation/pkg/models"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type FilePostgres struct {
	db *sqlx.DB
}

func NewFilePostgres(db *sqlx.DB) *FilePostgres {
	return &FilePostgres{db: db}
}

func (r *FilePostgres) Create(userId int, file models.File) (int, error) {
    var id int
    query := `INSERT INTO files (user_id, title, path, status, file_content) 
              VALUES ($1, $2, $3, $4, $5) RETURNING id`
    
    err := r.db.QueryRow(query, userId, file.Title, file.Path, file.Status, file.FileContent).Scan(&id)
    return id, err
}

func (r *FilePostgres) GetAll(userId int) ([]models.File, error) {
	var files []models.File
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1 ORDER BY created_at DESC", filesTable)
	err := r.db.Select(&files, query, userId)
	return files, err
}

func (r *FilePostgres) GetById(userId, fileId int) (models.File, error) {
	var file models.File
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1 AND id = $2", filesTable)
	err := r.db.Get(&file, query, userId, fileId)
	return file, err
}

func (r *FilePostgres) Delete(userId, fileId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1 AND id = $2", filesTable)
	_, err := r.db.Exec(query, userId, fileId)
	return err
}

func (r *FilePostgres) UpdateStatus(fileId int, status string) error {
	query := fmt.Sprintf("UPDATE %s SET status = $1 WHERE id = $2", filesTable)
	_, err := r.db.Exec(query, status, fileId)
	return err
}