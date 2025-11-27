package model

import (
	"time"
	
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Translation struct {
	ID            int       `json:"id"`
	FileID        int       `json:"file_id"`
	SourceLang    string    `json:"source_lang"`
	TargetLang    string    `json:"target_lang"`
	Status        string    `json:"status"` // pending, processing, completed, failed
	TranslatedPath string   `json:"translated_path"`
	Error         string    `json:"error"`
	CreatedAt     time.Time `json:"created_at"`
	CompletedAt   time.Time `json:"completed_at"`
}

func (t *Translation) Validate() error {
	return validation.ValidateStruct(
		t,
		validation.Field(&t.SourceLang, validation.Required, validation.Length(2, 10)),
		validation.Field(&t.TargetLang, validation.Required, validation.Length(2, 10)),
		validation.Field(&t.FileID, validation.Required),
	)
}