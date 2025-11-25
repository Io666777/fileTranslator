package apiserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Io666777/fileTranslator/internal/app/store/sqlstore/teststore"
	"github.com/stretchr/testify/assert"
)

func TestServer_HandleUsersCreate(t *testing.T) {
    s := newServer(teststore.New())
    
    rec := httptest.NewRecorder()
    req, _ := http.NewRequest(http.MethodPost, "/users", nil)
    
    s.ServeHTTP(rec, req)
    assert.Equal(t, http.StatusOK, rec.Code) // Правильный порядок аргументов
}