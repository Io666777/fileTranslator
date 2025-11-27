package store

import "github.com/Io666777/fileTranslator/internal/app/model"

type UserRepository interface {
	Create(*model.User) error
	Find(int) (*model.User, error)
	FindByEmail(string) (*model.User, error)
}

type FileRepository interface {
	Create(*model.File) error
	Find(int) (*model.File, error)
	FindByUserID(int) ([]*model.File, error)
	Delete(int) error
	Update(*model.File) error
}

type TranslationRepository interface {
	Create(*model.Translation) error
	Find(int) (*model.Translation, error)
	FindByFileID(int) (*model.Translation, error)
	FindByUserID(int) ([]*model.Translation, error)
	Update(*model.Translation) error
}