package service

import "filetranslation/pkg/repository"

type Authorization interface{}

type Fileb interface{}

type Filea interface{}

type Service struct {
	Authorization
	Fileb
	Filea
}

func NewService(repos *repository.Repository) *Service {
	return &Service{}
}
