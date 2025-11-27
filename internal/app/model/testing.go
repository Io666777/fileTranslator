package model

import "testing"

func TestUser(t *testing.T) *User{
	return  &User{
		Email: "user@example.org",
		Password: "password",
	}
}

func TestFile(t *testing.T) *File {
	return &File{
		Filename: "test.txt",
		FileSize: 1024,
		MimeType: "text/plain",
		Status:   "uploaded",
	}
}

func TestTranslation(t *testing.T) *Translation {
	return &Translation{
		SourceLang: "en",
		TargetLang: "ru",
		Status:     "pending",
	}
}