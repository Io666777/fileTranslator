package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

	if len(text) > 450 {
		return s.translateLongText(text, fromLang, toLang)
	}

	return s.translateShortText(text, fromLang, toLang)
}

func (s *TranslationService) translateShortText(text, fromLang, toLang string) (string, error) {
	encodedText := url.QueryEscape(text)
	apiURL := fmt.Sprintf("%s/get?q=%s&langpair=%s|%s",
		s.apiURL, encodedText, fromLang, toLang)

	logrus.Debugf("Calling MyMemory API for %d chars", len(text))

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(apiURL)
	if err != nil {
		return "", fmt.Errorf("connection error: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read error: %w", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return "", fmt.Errorf("parse error: %w", err)
	}

	if details, ok := data["responseDetails"].(string); ok && details != "" {
		return "", fmt.Errorf("API error: %s", details)
	}

	if responseData, ok := data["responseData"].(map[string]interface{}); ok {
		if translated, ok := responseData["translatedText"].(string); ok && translated != "" {
			logrus.Debugf("Translated %d -> %d chars", len(text), len(translated))
			return translated, nil
		}
	}

	return "", fmt.Errorf("no translation found")
}

func (s *TranslationService) translateLongText(text, fromLang, toLang string) (string, error) {
	logrus.Infof("Splitting long text (%d chars) into chunks", len(text))

	chunks := splitIntoChunks(text, 400)

	var translatedChunks []string

	for i, chunk := range chunks {
		logrus.Debugf("Translating chunk %d/%d (%d chars)", i+1, len(chunks), len(chunk))

		translated, err := s.translateShortText(chunk, fromLang, toLang)
		if err != nil {
			logrus.Warnf("Failed to translate chunk %d: %v", i+1, err)

			translated = chunk
		}

		translatedChunks = append(translatedChunks, translated)

		if i < len(chunks)-1 {
			time.Sleep(100 * time.Millisecond)
		}
	}

	result := strings.Join(translatedChunks, " ")
	logrus.Infof("Long translation completed: %d -> %d chars",
		len(text), len(result))

	return result, nil
}

func splitIntoChunks(text string, maxLen int) []string {
	var chunks []string

	sentences := splitSentences(text)

	var currentChunk strings.Builder
	for _, sentence := range sentences {

		if currentChunk.Len()+len(sentence)+1 > maxLen && currentChunk.Len() > 0 {
			chunks = append(chunks, currentChunk.String())
			currentChunk.Reset()
		}

		if currentChunk.Len() > 0 {
			currentChunk.WriteString(" ")
		}
		currentChunk.WriteString(sentence)
	}

	if currentChunk.Len() > 0 {
		chunks = append(chunks, currentChunk.String())
	}

	if len(chunks) == 0 || (len(chunks) == 1 && len(chunks[0]) > maxLen) {
		chunks = splitByWords(text, maxLen)
	}

	return chunks
}

func splitSentences(text string) []string {
	// Заменяем разные виды точек
	text = strings.ReplaceAll(text, "。", ".")
	text = strings.ReplaceAll(text, "！", "!")
	text = strings.ReplaceAll(text, "？", "?")

	splitPatterns := []string{". ", "! ", "? ", ".\n", "!\n", "?\n", ".", "!", "?", "\n"}

	for _, pattern := range splitPatterns {
		if strings.Contains(text, pattern) {
			var sentences []string
			parts := strings.Split(text, pattern)

			for i, part := range parts {
				part = strings.TrimSpace(part)
				if part == "" {
					continue
				}

				if i < len(parts)-1 {
					part += pattern[0:1] 
				}

				sentences = append(sentences, part)
			}

			if len(sentences) > 1 {
				return sentences
			}
		}
	}

	return []string{text}
}

func splitByWords(text string, maxLen int) []string {
	var chunks []string
	var current strings.Builder

	words := strings.Fields(text)
	for _, word := range words {
		if current.Len()+len(word)+1 > maxLen && current.Len() > 0 {
			chunks = append(chunks, current.String())
			current.Reset()
		}

		if current.Len() > 0 {
			current.WriteString(" ")
		}
		current.WriteString(word)
	}

	if current.Len() > 0 {
		chunks = append(chunks, current.String())
	}

	return chunks
}
