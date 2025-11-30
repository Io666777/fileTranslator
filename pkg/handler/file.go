package handler

import (
    "filetranslation/pkg/models"
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
)

func (h *Handler) createFile(c *gin.Context) {
    userId, err := getUserId(c)
    if err != nil {
        newErrorResponse(c, http.StatusInternalServerError, err.Error())
        return
    }

    var input models.File
    if err := c.BindJSON(&input); err != nil {
        newErrorResponse(c, http.StatusBadRequest, err.Error())
        return
    }

    // ИСПРАВЛЕНО: правильно обрабатываем возвращаемые значения
    id, err := h.services.File.Create(userId, input)
    if err != nil {
        newErrorResponse(c, http.StatusInternalServerError, err.Error())
        return
    }

    c.JSON(http.StatusOK, map[string]interface{}{
        "id": id,
    })
}

func (h *Handler) getAllFiles(c *gin.Context) {
    userId, err := getUserId(c)
    if err != nil {
        newErrorResponse(c, http.StatusInternalServerError, err.Error())
        return
    }

    // ИСПРАВЛЕНО: правильно обрабатываем возвращаемые значения
    files, err := h.services.File.GetAll(userId)
    if err != nil {
        newErrorResponse(c, http.StatusInternalServerError, err.Error())
        return
    }

    c.JSON(http.StatusOK, files)
}

func (h *Handler) getFileById(c *gin.Context) {
    userId, err := getUserId(c)
    if err != nil {
        newErrorResponse(c, http.StatusInternalServerError, err.Error())
        return
    }

    // ИСПРАВЛЕНО: strconv.Atoi возвращает (int, error)
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        newErrorResponse(c, http.StatusBadRequest, "invalid id param")
        return
    }

    // ИСПРАВЛЕНО: правильно обрабатываем возвращаемые значения
    file, err := h.services.File.GetById(userId, id)
    if err != nil {
        newErrorResponse(c, http.StatusInternalServerError, err.Error())
        return
    }

    c.JSON(http.StatusOK, file)
}

func (h *Handler) deleteFile(c *gin.Context) {
    userId, err := getUserId(c)
    if err != nil {
        newErrorResponse(c, http.StatusInternalServerError, err.Error())
        return
    }

    // ИСПРАВЛЕНО: strconv.Atoi возвращает (int, error)
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        newErrorResponse(c, http.StatusBadRequest, "invalid id param")
        return
    }

    // ИСПРАВЛЕНО: правильно обрабатываем возвращаемые значения
    err = h.services.File.Delete(userId, id)
    if err != nil {
        newErrorResponse(c, http.StatusInternalServerError, err.Error())
        return
    }

    c.JSON(http.StatusOK, map[string]interface{}{
        "status": "ok",
    })
}

// createTranslation остаётся пока заглушкой
func (h *Handler) createTranslation(c *gin.Context) {
    c.JSON(http.StatusOK, map[string]interface{}{
        "message": "createTranslation not implemented yet",
    })
}