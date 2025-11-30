package service

import (
	"filetranslation/pkg/models"
	"filetranslation/pkg/repository"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
}

type Fileb interface{}

type Filea interface{}

type Service struct {
	Authorization
	Fileb
	Filea
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
	}
}
