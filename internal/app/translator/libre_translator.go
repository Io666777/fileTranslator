// internal/app/translator/libre_translator.go
package translator

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

type LibreTranslator struct {
    baseURL    string
    client     *http.Client
    // apiKey     string // опционально
}

type TranslateRequest struct {
    Q      string `json:"q"`
    Source string `json:"source"`
    Target string `json:"target"`
    Format string `json:"format"`
}

type TranslateResponse struct {
    TranslatedText string `json:"translatedText"`
}

func NewLibreTranslator(baseURL string) *LibreTranslator {
    return &LibreTranslator{
        baseURL: baseURL,
        client: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}

func (l *LibreTranslator) Translate(text, sourceLang, targetLang string) (string, error) {
    if text == "" {
        return "", nil
    }

    requestBody := TranslateRequest{
        Q:      text,
        Source: sourceLang,
        Target: targetLang,
        Format: "text",
    }

    jsonData, err := json.Marshal(requestBody)
    if err != nil {
        return "", fmt.Errorf("failed to marshal request: %w", err)
    }

    resp, err := l.client.Post(l.baseURL+"/translate", "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return "", fmt.Errorf("translation request failed: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return "", fmt.Errorf("translation failed: %s - %s", resp.Status, string(body))
    }

    var result TranslateResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return "", fmt.Errorf("failed to decode response: %w", err)
    }

    return result.TranslatedText, nil
}

func (l *LibreTranslator) SupportsFileType(mimeType string) bool {
    // Поддерживаем текстовые файлы
    supportedTypes := map[string]bool{
        "text/plain": true,
        "text/html":  true,
        "text/xml":   true,
        "application/json": true,
    }
    return supportedTypes[mimeType]
}

func (l *LibreTranslator) TranslateFile(filePath, sourceLang, targetLang string) (string, error) {
    // Реализуем позже - нужно читать файл и обрабатывать
    return "", fmt.Errorf("file translation not implemented yet")
}