package repository

import (
	"filetranslation/pkg/models"

	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
}

type Fileb interface{}

type Filea interface{}

type Repository struct {
	Authorization
	Fileb
	Filea
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
