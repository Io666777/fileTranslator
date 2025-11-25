package sqlstore

import (
	"database/sql"

	"github.com/Io666777/fileTranslator/internal/app/store"
	_ "github.com/lib/pq"
)

type Store struct {

	db             *sql.DB
	userRepository *UserRepository
}
//New
func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

//user
func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}
	return s.userRepository
}
