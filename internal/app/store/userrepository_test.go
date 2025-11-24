package store_test

import (
	"testing"

	"github.com/Io666777/fileTranslator/internal/app/model"
	"github.com/Io666777/fileTranslator/internal/app/store"

	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	databaseURL := "postgres://postgres:password@localhost:5432/ft_test?sslmode=disable"
	s, teardown := store.NewTestStore(t, databaseURL)
	defer teardown("users")

	u, err := s.User().Create(&model.User{
		Email:             "user@example.org",
		EncryptedPassword: "password",
	})
	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	databaseURL := "postgres://postgres:password@localhost:5432/ft_test?sslmode=disable" // Добавьте эту строку
	s, teardown := store.NewTestStore(t, databaseURL)
	defer teardown("users")

	email := "user@example.org"
	_, err := s.User().FindByEmail(email)
	assert.Error(t, err)

	s.User().Create(&model.User{
		Email:             "user@example.org",
		EncryptedPassword: "password",
	})
	u, err := s.User().FindByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}
