package sqlstore

import (
	"database/sql"

	"github.com/Io666777/fileTranslator/internal/app/store"
	_ "github.com/lib/pq"
)

type Store struct {
	db                    *sql.DB
	userRepository        *UserRepository
	fileRepository        *FileRepository
	translationRepository *TranslationRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}
	return s.userRepository
}

func (s *Store) File() store.FileRepository {
	if s.fileRepository != nil {
		return s.fileRepository
	}

	s.fileRepository = &FileRepository{
		store: s,
	}
	return s.fileRepository
}

func (s *Store) Translation() store.TranslationRepository {
	if s.translationRepository != nil {
		return s.translationRepository
	}

	s.translationRepository = &TranslationRepository{
		store: s,
	}
	return s.translationRepository
}