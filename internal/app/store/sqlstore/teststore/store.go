package teststore

import (
	"github.com/Io666777/fileTranslator/internal/app/model"
	"github.com/Io666777/fileTranslator/internal/app/store"
)

type Store struct {
	userRepository *UserRepository
}
//New
func New() *Store {
	return &Store{}
}

//user
func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
		users: make(map[string]*model.User),
	}
	return s.userRepository
}
