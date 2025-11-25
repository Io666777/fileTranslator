package model_test

import (
	"testing"

	"github.com/Io666777/fileTranslator/internal/app/model"
	"github.com/stretchr/testify/assert"
)

func TestUser_BeforeCreate(t *testing.T) {
	u := model.TestUser(t)
	assert.NoError(t,u.BeforeCreate())
	assert.NotEmpty(t,u.EncryptedPassword)
}