package handler

import (
	"bytes"
	"filetranslation/pkg/models"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) uploadFile(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "file is required")
		return
	}

	logrus.Infof("Uploading file: %s, size: %d, user: %d",
		file.Filename, file.Size, userId)

	src, err := file.Open()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "failed to open file")
		return
	}
	defer src.Close()

	var buf bytes.Buffer
	written, err := io.Copy(&buf, src)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "failed to read file: "+err.Error())
		return
	}

	fileContent := buf.Bytes()

	logrus.Infof("File read successfully: expected %d bytes, read %d bytes",
		file.Size, written)

	fileRecord := models.File{
		Title:       file.Filename,
		Path:        "db",
		Status:      "uploaded",
		UserID:      userId,
		FileContent: fileContent,
	}

	id, err := h.services.File.Create(userId, fileRecord)
	if err != nil {
		logrus.Errorf("Failed to save file to DB: %v", err)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	logrus.Infof("File uploaded successfully: id=%d, title=%s, content length=%d",
		id, file.Filename, len(fileContent))

	c.JSON(http.StatusOK, map[string]interface{}{"id": id})
}

func (h *Handler) getAllFiles(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	logrus.Infof("Getting files for user: %d", userId)

	files, err := h.services.File.GetAll(userId)
	if err != nil {
		logrus.Errorf("Failed to get files: %v", err)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	logrus.Infof("Retrieved %d files for user %d", len(files), userId)

	for i, file := range files {
		logrus.Debugf("File %d: id=%d, title=%s, status=%s, content length=%d",
			i, file.ID, file.Title, file.Status, len(file.FileContent))
	}

	c.JSON(http.StatusOK, files)
}

func (h *Handler) downloadFile(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	logrus.Infof("Downloading file: id=%d, user=%d", id, userId)

	file, err := h.services.File.GetById(userId, id)
	if err != nil {
		logrus.Errorf("File not found: %v", err)
		newErrorResponse(c, http.StatusInternalServerError, "file not found")
		return
	}

	logrus.Infof("File found: id=%d, title=%s, content length=%d",
		file.ID, file.Title, len(file.FileContent))

	c.Header("Content-Disposition", "attachment; filename="+file.Title)
	c.Data(http.StatusOK, "application/octet-stream", file.FileContent)
}

func (h *Handler) createTranslation(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	logrus.Infof("Starting translation for file %d, user %d", id, userId)

	file, err := h.services.File.GetById(userId, id)
	if err != nil {
		logrus.Errorf("File not found: %v", err)
		newErrorResponse(c, http.StatusInternalServerError, "file not found")
		return
	}

	logrus.Infof("Original file: id=%d, title=%s, content length=%d bytes",
		file.ID, file.Title, len(file.FileContent))

	contentText := string(file.FileContent)

	fromLang, toLang := detectLanguage(contentText)

	logrus.Infof("Language detection: %s -> %s", fromLang, toLang)

	contentPreview := contentText
	if len(contentPreview) > 100 {
		contentPreview = contentPreview[:100] + "..."
	}
	logrus.Infof("File content preview: %s", contentPreview)

	err = h.services.File.UpdateStatus(id, "processing")
	if err != nil {
		logrus.Errorf("Failed to update status: %v", err)
	}

	translatedText, err := h.services.Translation.TranslateText(
		contentText,
		fromLang,
		toLang,
	)

	if err != nil {
		logrus.Errorf("Translation error: %v", err)
		h.services.File.UpdateStatus(id, "error")
		newErrorResponse(c, http.StatusInternalServerError, "translation failed: "+err.Error())
		return
	}

	logrus.Infof("Translation successful: %d -> %d bytes",
		len(contentText), len(translatedText))

	var translatedTitle string
	if fromLang == "ru" && toLang == "en" {
		translatedTitle = "translated_en_" + file.Title
	} else if fromLang == "en" && toLang == "ru" {
		translatedTitle = "translated_ru_" + file.Title
	} else {
		translatedTitle = "translated_" + file.Title
	}

	translatedFile := models.File{
		Title:       translatedTitle,
		Path:        "db",
		Status:      "translated",
		UserID:      userId,
		FileContent: []byte(translatedText),
	}

	translatedId, err := h.services.File.Create(userId, translatedFile)
	if err != nil {
		logrus.Errorf("Failed to save translation: %v", err)
		newErrorResponse(c, http.StatusInternalServerError, "failed to save translation")
		return
	}

	h.services.File.UpdateStatus(id, "completed")

	logrus.Infof("Translation completed: %s->%s, new file ID: %d",
		fromLang, toLang, translatedId)

	c.JSON(http.StatusOK, map[string]interface{}{
		"message":            fmt.Sprintf("translation completed (%s->%s)", fromLang, toLang),
		"original_file_id":   id,
		"translated_file_id": translatedId,
		"from_lang":          fromLang,
		"to_lang":            toLang,
	})
}

func (h *Handler) deleteFile(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	logrus.Infof("Deleting file: id=%d, user=%d", id, userId)

	err = h.services.File.Delete(userId, id)
	if err != nil {
		logrus.Errorf("Failed to delete file: %v", err)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	logrus.Infof("File deleted: id=%d", id)

	c.JSON(http.StatusOK, map[string]interface{}{"status": "deleted"})
}

func detectLanguage(text string) (string, string) {
	if text == "" {
		return "ru", "en"
	}

	var russianChars, englishChars int
	var totalChars int

	for _, char := range text {
		if unicode.IsLetter(char) {
			totalChars++

			if (char >= 'а' && char <= 'я') || (char >= 'А' && char <= 'Я') || char == 'ё' || char == 'Ё' {
				russianChars++
			}

			if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
				englishChars++
			}
		}
	}

	if totalChars == 0 {
		return "ru", "en"
	}

	russianPercent := float64(russianChars) * 100 / float64(totalChars)
	englishPercent := float64(englishChars) * 100 / float64(totalChars)

	logrus.Debugf("Language detection: Russian %.1f%%, English %.1f%%",
		russianPercent, englishPercent)

	if russianPercent > 50 {
		return "ru", "en"
	} else if englishPercent > 50 {
		return "en", "ru"
	} else {

		textLower := strings.ToLower(text)

		russianKeywords := []string{"привет", "мир", "спасибо", "пожалуйста", "да", "нет"}
		englishKeywords := []string{"hello", "world", "thank", "please", "yes", "no"}

		russianMatches := 0
		englishMatches := 0

		for _, keyword := range russianKeywords {
			if strings.Contains(textLower, keyword) {
				russianMatches++
			}
		}

		for _, keyword := range englishKeywords {
			if strings.Contains(textLower, keyword) {
				englishMatches++
			}
		}

		if russianMatches > englishMatches {
			return "ru", "en"
		} else if englishMatches > russianMatches {
			return "en", "ru"
		}
	}

	return "ru", "en"
}
