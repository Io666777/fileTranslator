package models // исправлено filetranslator -> models

type User struct {
	ID       int    `json:"-" gorm:"primaryKey"`
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" gorm:"unique" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type File struct { // объединяем Filebt и Fileat в одну структуру
	ID    int    `json:"id" gorm:"primaryKey"`
	Title string `json:"title"`
	Path  string `json:"path"`
	// Добавляем поля для перевода
	OriginalText   string `json:"original_text"`
	TranslatedText string `json:"translated_text"`
	SourceLanguage string `json:"source_language"`
	TargetLanguage string `json:"target_language"`
	Status         string `json:"status"` // uploaded, processing, translated
}
