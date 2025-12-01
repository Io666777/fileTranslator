package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type TranslationService struct {
	apiURL string
}

func NewTranslationService(apiURL string) *TranslationService {
	return &TranslationService{apiURL: apiURL}
}

func (s *TranslationService) TranslateText(text, fromLang, toLang string) (string, error) {
	if text == "" {
		return "", nil
	}

	logrus.Infof("Translating %d chars from %s to %s", len(text), fromLang, toLang)

	// Очищаем текст
	text = strings.TrimSpace(text)
	if text == "" {
		return "", nil
	}

	// Формируем запрос для LibreTranslate
	request := struct {
		Q      string `json:"q"`
		Source string `json:"source"`
		Target string `json:"target"`
		Format string `json:"format"`
	}{
		Q:      text,
		Source: fromLang,
		Target: toLang,
		Format: "text",
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("marshal error: %w", err)
	}

	// URL для LibreTranslate
	url := strings.TrimSuffix(s.apiURL, "/") + "/translate"

	// Логируем что отправляем
	if len(text) < 100 {
		logrus.Debugf("Sending to translate: '%s'", text)
	} else {
		logrus.Debugf("Sending to translate: '%s...' (%d chars total)", text[:100], len(text))
	}

	// Создаем запрос
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("request error: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Отправляем запрос
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("connection error: %w", err)
	}
	defer resp.Body.Close()

	// Читаем ответ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read error: %w", err)
	}

	// Логируем сырой ответ
	logrus.Debugf("Translation response status: %d", resp.StatusCode)
	logrus.Debugf("Translation response body: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("translation failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Парсим ответ
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("parse error: %w, body: %s", err, string(body))
	}

	// Извлекаем переведенный текст
	translatedText, ok := result["translatedText"].(string)
	if !ok {
		return "", fmt.Errorf("no translatedText in response: %v", result)
	}

	if translatedText == "" {
		return "", fmt.Errorf("empty translation received")
	}

	logrus.Infof("Translation successful: %d -> %d chars", 
		len(text), len(translatedText))
	
	// Логируем результат
	if len(translatedText) < 100 {
		logrus.Debugf("Translated: '%s'", translatedText)
	} else {
		logrus.Debugf("Translated: '%s...'", translatedText[:100])
	}

	return translatedText, nil
}