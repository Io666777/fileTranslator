package teststore_test

import (
	"testing"

	"github.com/Io666777/fileTranslator/internal/app/model"
	"github.com/Io666777/fileTranslator/internal/app/store"
	"github.com/Io666777/fileTranslator/internal/app/store/sqlstore/teststore"

	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	 
	s := teststore.New()
	u := model.TestUser(t)
	 
	err := s.User().Create(u) // Теперь получаем только error
	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	
	s := teststore.New()
	email := "user@example.org"

	// Сначала проверяем, что пользователя нет
	_, err := s.User().FindByEmail(email)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	// Создаем пользователя
	testUser := model.TestUser(t)
	testUser.Email = email
	err = s.User().Create(testUser) // ← Исправлено: убрал лишнюю переменную
	assert.NoError(t, err)

	// Ищем созданного пользователя
	foundUser, err := s.User().FindByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
}