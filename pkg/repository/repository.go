package repository

import (
	"filetranslation/pkg/models"

	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GetUser(username, password string) (models.User, error)
}

type File interface {
	Create(userId int, file models.File) (int, error)
	GetAll(userId int) ([]models.File, error)
	GetById(userId, fileId int) (models.File, error)
	Delete(userId, fileId int) error
	UpdateStatus(fileId int, status string) error
}

type Repository struct {
	Authorization
	File
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		File:          NewFilePostgres(db), //cannot use NewAuthPostgres(db) (value of type *AuthPostgres) as Authorization value in struct literal: *AuthPostgres does not implement Authorization (wrong type for method CreateUser)have CreateUser(string, string) (models.User, error)want CreateUser(models.User) (int, error)
	}
}
