package handler

import (
	"filetranslation/pkg/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

	// Читаем файл
	src, err := file.Open()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "failed to open file")
		return
	}
	defer src.Close()

	fileContent := make([]byte, file.Size)
	_, err = src.Read(fileContent) // ДОБАВИТЬ ПРОВЕРКУ ОШИБКИ
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "failed to read file")
		return
	}
	src.Read(fileContent)

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
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"id": id})
}

// СПИСОК ВСЕХ ФАЙЛОВ
func (h *Handler) getAllFiles(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	files, err := h.services.File.GetAll(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
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

	file, err := h.services.File.GetById(userId, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "file not found")
		return
	}

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

	userId, _ := getUserId(c)

	// Получаем файл из БД
	file, err := h.services.File.GetById(userId, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "file not found")
		return
	}

	// Обновляем статус
	h.services.File.UpdateStatus(id, "processing")

	// ПРОСТЕЙШАЯ логика перевода (только текст)
	translatedText, err := h.services.Translation.TranslateText(
		string(file.FileContent),
		"auto",
		"en", // переводим на английский
	)

	if err != nil {
		h.services.File.UpdateStatus(id, "error")
		newErrorResponse(c, http.StatusInternalServerError, "translation failed")
		return
	}

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
		newErrorResponse(c, http.StatusInternalServerError, "failed to save translation")
		return
	}

	h.services.File.UpdateStatus(id, "completed")

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

	err = h.services.File.Delete(userId, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"status": "deleted"})
}
