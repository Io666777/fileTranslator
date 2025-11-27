package teststore

import (
	"github.com/Io666777/fileTranslator/internal/app/model"
	"github.com/Io666777/fileTranslator/internal/app/store"
)

type FileRepository struct {
	store *Store
	files map[int]*model.File
}

func (r *FileRepository) Create(f *model.File) error {
	if err := f.Validate(); err != nil {
		return err
	}

	f.ID = len(r.files) + 1
	r.files[f.ID] = f
	return nil
}

func (r *FileRepository) Find(id int) (*model.File, error) {
	f, ok := r.files[id]
	if !ok {
		return nil, store.ErrRecordNotFound
	}
	return f, nil
}

func (r *FileRepository) FindByUserID(userID int) ([]*model.File, error) {
	var files []*model.File
	for _, file := range r.files {
		if file.UserID == userID {
			files = append(files, file)
		}
	}
	return files, nil
}

func (r *FileRepository) Delete(id int) error {
	delete(r.files, id)
	return nil
}

func (r *FileRepository) Update(f *model.File) error {
	r.files[f.ID] = f
	return nil
}