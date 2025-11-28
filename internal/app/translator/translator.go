// internal/app/translator/translator.go
package translator

type Translator interface {
    Translate(text, sourceLang, targetLang string) (string, error)
    SupportsFileType(mimeType string) bool
    TranslateFile(filePath, sourceLang, targetLang string) (string, error)
}