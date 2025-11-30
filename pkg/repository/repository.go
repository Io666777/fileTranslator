package repository

import (
	"filetranslation/pkg/models"

	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GetUser(username, password string) (models.User, error)
}



type Repository struct {
	Authorization

}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),//cannot use NewAuthPostgres(db) (value of type *AuthPostgres) as Authorization value in struct literal: *AuthPostgres does not implement Authorization (wrong type for method CreateUser)have CreateUser(string, string) (models.User, error)want CreateUser(models.User) (int, error)
	}
}
