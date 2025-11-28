// internal/app/translator/file_processor.go
package translator

import (
    "io"
    "os"
    "strings"
)

type FileProcessor struct{}

func NewFileProcessor() *FileProcessor {
    return &FileProcessor{}
}

// ExtractText читает текст из файла (пока только txt файлы)
func (fp *FileProcessor) ExtractText(filePath string) (string, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return "", err
    }
    defer file.Close()

    content, err := io.ReadAll(file)
    if err != nil {
        return "", err
    }

    return string(content), nil
}

// SaveTranslatedText сохраняет переведенный текст в файл
func (fp *FileProcessor) SaveTranslatedText(text, filePath string) (string, error) {
    // Создаем путь для переведенного файла
    translatedPath := strings.TrimSuffix(filePath, ".txt") + "_translated.txt"
    
    file, err := os.Create(translatedPath)
    if err != nil {
        return "", err
    }
    defer file.Close()

    _, err = file.WriteString(text)
    if err != nil {
        return "", err
    }
    
    return translatedPath, nil
}

// SupportsFile проверяет поддерживаемые форматы
func (fp *FileProcessor) SupportsFile(mimeType string) bool {
    supported := map[string]bool{
        "text/plain": true,
        "text/html":  true,
    }
    return supported[mimeType]
}