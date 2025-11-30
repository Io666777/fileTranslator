package service

import (
    "bytes"
    "encoding/json"
    "net/http"
)

type TranslationService struct {
    apiURL string // "http://localhost:5000"
}

func NewTranslationService(apiURL string) *TranslationService {
    return &TranslationService{apiURL: apiURL}
}

// Простейший метод перевода текста
func (s *TranslationService) TranslateText(text, fromLang, toLang string) (string, error) {
    requestBody := map[string]string{
        "q":      text,
        "source": fromLang,
        "target": toLang,
        "format": "text",
    }
    
    jsonData, _ := json.Marshal(requestBody)
    
    resp, err := http.Post(s.apiURL+"/translate", "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    
    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)
    
    translatedText, _ := result["translatedText"].(string)
    return translatedText, nil
}