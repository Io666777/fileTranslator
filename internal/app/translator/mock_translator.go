// internal/app/translator/mock_translator.go
package translator

// Мок для тестирования без сервера
type MockTranslator struct{}

func NewMockTranslator() *MockTranslator {
    return &MockTranslator{}
}

func (m *MockTranslator) Translate(text, sourceLang, targetLang string) (string, error) {
    if text == "" {
        return "", nil
    }
    return "[TRANSLATED] " + text, nil
}

func (m *MockTranslator) SupportsFileType(mimeType string) bool {
    return mimeType == "text/plain"
}

func (m *MockTranslator) TranslateFile(filePath, sourceLang, targetLang string) (string, error) {
    return filePath + ".translated", nil
}