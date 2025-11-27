package teststore

import (
	"github.com/Io666777/fileTranslator/internal/app/model"
	"github.com/Io666777/fileTranslator/internal/app/store"
)

type TranslationRepository struct {
	store        *Store
	translations map[int]*model.Translation
}

func (r *TranslationRepository) Create(t *model.Translation) error {
	if err := t.Validate(); err != nil {
		return err
	}

	t.ID = len(r.translations) + 1
	r.translations[t.ID] = t
	return nil
}

func (r *TranslationRepository) Find(id int) (*model.Translation, error) {
	t, ok := r.translations[id]
	if !ok {
		return nil, store.ErrRecordNotFound
	}
	return t, nil
}

func (r *TranslationRepository) FindByFileID(fileID int) (*model.Translation, error) {
	for _, t := range r.translations {
		if t.FileID == fileID {
			return t, nil
		}
	}
	return nil, store.ErrRecordNotFound
}

func (r *TranslationRepository) FindByUserID(userID int) ([]*model.Translation, error) {
	var translations []*model.Translation
	for _, t := range r.translations {
		file, err := r.store.File().Find(t.FileID)
		if err == nil && file.UserID == userID {
			translations = append(translations, t)
		}
	}
	return translations, nil
}

func (r *TranslationRepository) Update(t *model.Translation) error {
	r.translations[t.ID] = t
	return nil
}