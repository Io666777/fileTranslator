package sqlstore

import (
	"database/sql"

	"github.com/Io666777/fileTranslator/internal/app/model"
	"github.com/Io666777/fileTranslator/internal/app/store"
)

type TranslationRepository struct {
	store *Store
}

func (r *TranslationRepository) Create(t *model.Translation) error {
	if err := t.Validate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO translations (file_id, source_lang, target_lang, status) VALUES ($1, $2, $3, $4) RETURNING id",
		t.FileID,
		t.SourceLang,
		t.TargetLang,
		t.Status,
	).Scan(&t.ID)
}

func (r *TranslationRepository) Find(id int) (*model.Translation, error) {
	t := &model.Translation{}
	err := r.store.db.QueryRow(
		"SELECT id, file_id, source_lang, target_lang, status, translated_path, error, created_at, completed_at FROM translations WHERE id = $1",
		id,
	).Scan(&t.ID, &t.FileID, &t.SourceLang, &t.TargetLang, &t.Status, &t.TranslatedPath, &t.Error, &t.CreatedAt, &t.CompletedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	return t, nil
}

func (r *TranslationRepository) FindByFileID(fileID int) (*model.Translation, error) {
	t := &model.Translation{}
	err := r.store.db.QueryRow(
		"SELECT id, file_id, source_lang, target_lang, status, translated_path, error, created_at, completed_at FROM translations WHERE file_id = $1",
		fileID,
	).Scan(&t.ID, &t.FileID, &t.SourceLang, &t.TargetLang, &t.Status, &t.TranslatedPath, &t.Error, &t.CreatedAt, &t.CompletedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	return t, nil
}

func (r *TranslationRepository) FindByUserID(userID int) ([]*model.Translation, error) {
	rows, err := r.store.db.Query(
		`SELECT t.id, t.file_id, t.source_lang, t.target_lang, t.status, t.translated_path, t.error, t.created_at, t.completed_at 
		 FROM translations t 
		 JOIN files f ON t.file_id = f.id 
		 WHERE f.user_id = $1 
		 ORDER BY t.created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var translations []*model.Translation
	for rows.Next() {
		t := &model.Translation{}
		err := rows.Scan(&t.ID, &t.FileID, &t.SourceLang, &t.TargetLang, &t.Status, &t.TranslatedPath, &t.Error, &t.CreatedAt, &t.CompletedAt)
		if err != nil {
			return nil, err
		}
		translations = append(translations, t)
	}

	return translations, nil
}

func (r *TranslationRepository) Update(t *model.Translation) error {
	_, err := r.store.db.Exec(
		"UPDATE translations SET status = $1, translated_path = $2, error = $3, completed_at = $4 WHERE id = $5",
		t.Status, t.TranslatedPath, t.Error, t.CompletedAt, t.ID,
	)
	return err
}