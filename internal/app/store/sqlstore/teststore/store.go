package teststore

import (
	"github.com/Io666777/fileTranslator/internal/app/model"
	"github.com/Io666777/fileTranslator/internal/app/store"
)

type Store struct {
	userRepository        *UserRepository
	fileRepository        *FileRepository
	translationRepository *TranslationRepository
}

func New() *Store {
	return &Store{}
}

func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
		users: make(map[int]*model.User),
	}
	return s.userRepository
}

func (s *Store) File() store.FileRepository {
	if s.fileRepository != nil {
		return s.fileRepository
	}

	s.fileRepository = &FileRepository{
		store: s,
		files: make(map[int]*model.File),
	}
	return s.fileRepository
}

func (s *Store) Translation() store.TranslationRepository {
	if s.translationRepository != nil {
		return s.translationRepository
	}

	s.translationRepository = &TranslationRepository{
		store:        s,
		translations: make(map[int]*model.Translation),
	}
	return s.translationRepository
}