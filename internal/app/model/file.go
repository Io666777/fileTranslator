package model

import (
	"time"
	
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type File struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	Filename     string    `json:"filename"`
	OriginalPath string    `json:"original_path"`
	FileSize     int64     `json:"file_size"`
	MimeType     string    `json:"mime_type"`
	Status       string    `json:"status"` // uploaded, processing, translated, error
	CreatedAt    time.Time `json:"created_at"`
}

func (f *File) Validate() error {
	return validation.ValidateStruct(
		f,
		validation.Field(&f.Filename, validation.Required, validation.Length(1, 255)),
		validation.Field(&f.MimeType, validation.Required),
		validation.Field(&f.FileSize, validation.Required, validation.Min(1)),
	)
}