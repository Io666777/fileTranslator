package store

type Store interface {
	User() UserRepository
	File() FileRepository
	Translation() TranslationRepository
}