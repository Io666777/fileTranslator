package service

import (
	"filetranslation/pkg/models"
	"filetranslation/pkg/repository"
)

type FileService struct {
	repo repository.File
}

func NewFileService(repo repository.File) *FileService {
	return &FileService{repo: repo}
}

func (s *FileService) Create(userId int, file models.File) (int, error) {
	return s.repo.Create(userId, file)
}

func (s *FileService) GetAll(userId int) ([]models.File, error) {
	return s.repo.GetAll(userId)
}

func (s *FileService) GetById(userId, fileId int) (models.File, error) {
	return s.repo.GetById(userId, fileId)
}

func (s *FileService) Delete(userId, fileId int) error {
	return s.repo.Delete(userId, fileId)
}

func (s *FileService) UpdateStatus(fileId int, status string) error {
	return s.repo.UpdateStatus(fileId, status)
}