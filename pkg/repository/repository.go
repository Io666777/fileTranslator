package repository

type Authorization interface{}

type Fileb interface{}

type Filea interface{}

type Repository struct{
	Authorization
	Fileb
	Filea
}

func NewRepository() *Repository {
	return &Repository{}
}


