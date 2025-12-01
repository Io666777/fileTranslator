package handler

import (
	"bytes"
	"filetranslation/pkg/models"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ЗАГРУЗКА ФАЙЛА
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

	// Читаем файл
	src, err := file.Open()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "failed to open file")
		return
	}
	defer src.Close()

	// СПОСОБ 1: Используем bytes.Buffer для чтения
	var buf bytes.Buffer
	written, err := io.Copy(&buf, src)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "failed to read file: "+err.Error())
		return
	}

	fileContent := buf.Bytes()

	logrus.Infof("File read successfully: expected %d bytes, read %d bytes",
		file.Size, written)

	// Сохраняем в БД
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

// СПИСОК ВСЕХ ФАЙЛОВ
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

	// Логируем информацию о каждом файле
	for i, file := range files {
		logrus.Debugf("File %d: id=%d, title=%s, status=%s, content length=%d",
			i, file.ID, file.Title, file.Status, len(file.FileContent))
	}

	c.JSON(http.StatusOK, files)
}

// СКАЧИВАНИЕ ФАЙЛА
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

// ЗАПРОС ПЕРЕВОДА
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

	// Получаем файл из БД
	file, err := h.services.File.GetById(userId, id)
	if err != nil {
		logrus.Errorf("File not found: %v", err)
		newErrorResponse(c, http.StatusInternalServerError, "file not found")
		return
	}

	logrus.Infof("Original file: id=%d, title=%s, content length=%d bytes",
		file.ID, file.Title, len(file.FileContent))

	// Логируем первые 100 символов содержимого
	contentPreview := string(file.FileContent)
	if len(contentPreview) > 100 {
		contentPreview = contentPreview[:100] + "..."
	}
	logrus.Infof("File content preview: %s", contentPreview)

	// Обновляем статус
	err = h.services.File.UpdateStatus(id, "processing")
	if err != nil {
		logrus.Errorf("Failed to update status: %v", err)
	}

	// Переводим
	translatedText, err := h.services.Translation.TranslateText(
		string(file.FileContent),
		"auto",
		"en",
	)

	if err != nil {
		logrus.Errorf("Translation error: %v", err)
		h.services.File.UpdateStatus(id, "error")
		newErrorResponse(c, http.StatusInternalServerError, "translation failed: "+err.Error())
		return
	}

	logrus.Infof("Translation successful, translated length: %d bytes", len(translatedText))

	// Сохраняем переведенный файл
	translatedFile := models.File{
		Title:       "translated_" + file.Title,
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

	logrus.Infof("Translation completed, new file ID: %d, content length: %d",
		translatedId, len(translatedText))

	c.JSON(http.StatusOK, map[string]interface{}{
		"message":            "translation completed",
		"original_file_id":   id,
		"translated_file_id": translatedId,
	})
}

// УДАЛЕНИЕ ФАЙЛА
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
