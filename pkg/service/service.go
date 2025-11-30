package service

import (
	"filetranslation/pkg/models"
	"filetranslation/pkg/repository"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type File interface {
	Create(userId int, file models.File) (int, error)
	GetAll(userId int) ([]models.File, error)
	GetById(userId, fileId int) (models.File, error)
	Delete(userId, fileId int) error
	UpdateStatus(fileId int, status string) error
}

type Service struct {
	Authorization
	File
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		File:          NewFileService(repos.File),
	}
}
