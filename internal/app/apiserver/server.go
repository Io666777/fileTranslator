package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Io666777/fileTranslator/internal/app/model"
	"github.com/Io666777/fileTranslator/internal/app/store"
	"github.com/Io666777/fileTranslator/internal/app/translator"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

const (
	sessionName        = "ft"
	ctxKeyUser  ctxKey = iota
	ctxKeyRequestID
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated         = errors.New("not authenticated")
)

type ctxKey int8

type server struct {
	router        *mux.Router
	logger        *logrus.Logger
	store         store.Store
	sessionStore  sessions.Store
	translator    translator.Translator     // ДОБАВИТЬ
	fileProcessor *translator.FileProcessor // ДОБАВИТЬ
}

func newServer(store store.Store, sessionStore sessions.Store) *server {
	s := &server{
		router:        mux.NewRouter(),
		logger:        logrus.New(),
		store:         store,
		sessionStore:  sessionStore,
		translator:    translator.NewLibreTranslator("http://localhost:5000"), // ДОБАВИТЬ
		fileProcessor: translator.NewFileProcessor(),                          // ДОБАВИТЬ
	}
	s.configureRouter()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	s.router.HandleFunc("/users", s.handleUsersCreate()).Methods("POST")
	s.router.HandleFunc("/sessions", s.handleSessionsCreate()).Methods("POST")

	private := s.router.PathPrefix("/private").Subrouter()
	private.Use(s.authenticateUser)
	private.HandleFunc("/whoami", s.handleWhoami()).Methods("GET")

	private.HandleFunc("/files", s.handleFilesCreate()).Methods("POST")
	private.HandleFunc("/files", s.handleFilesList()).Methods("GET")
	private.HandleFunc("/files/{id:[0-9]+}", s.handleFilesGet()).Methods("GET")
	private.HandleFunc("/files/{id:[0-9]+}", s.handleFilesDelete()).Methods("DELETE")
	private.HandleFunc("/files/{id:[0-9]+}/translate", s.handleFilesTranslate()).Methods("POST")
	private.HandleFunc("/files/{id:[0-9]+}/download", s.handleFilesDownload()).Methods("GET")

	private.HandleFunc("/translations", s.handleTranslationsList()).Methods("GET")
	private.HandleFunc("/translations/{id:[0-9]+}", s.handleTranslationsGet()).Methods("GET")
}

func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxKeyRequestID),
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &ResponseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		logger.WithFields(logrus.Fields{
			"status_code": rw.code,
			"status":      http.StatusText(rw.code),
			"duration":    time.Since(start),
		}).Info("completed request")
	})
}
func (s *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		id, ok := session.Values["user_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		u, err := s.store.User().Find(id.(int))
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, u)))
	})
}

func (s *server) handleWhoami() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(ctxKeyUser).(*model.User)
		if !ok {
			s.error(w, r, http.StatusInternalServerError, errors.New("user not found in context"))
			return
		}
		s.respond(w, r, http.StatusOK, user)
	}
}

func (s *server) handleUsersCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u := &model.User{
			Email:    req.Email,
			Password: req.Password,
		}

		if err := s.store.User().Create(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		u.Sanitize()
		s.respond(w, r, http.StatusCreated, u)
	}
}

func (s *server) handleSessionsCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().FindByEmail(req.Email)
		if err != nil || !u.ComparePassword(req.Password) {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}

		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		session.Values["user_id"] = u.ID
		if err := s.sessionStore.Save(r, w, session); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func (s *server) handleFilesCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if err := r.ParseMultipartForm(10 << 20); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		defer file.Close()

		user := r.Context().Value(ctxKeyUser).(*model.User)

		filename := header.Filename
		filePath := fmt.Sprintf("storage/uploads/%d_%d_%s", user.ID, time.Now().Unix(), filename)

		os.MkdirAll("storage/uploads", 0755)

		dst, err := os.Create(filePath)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		defer dst.Close()

		fileSize, err := io.Copy(dst, file)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		f := &model.File{
			UserID:       user.ID,
			Filename:     filename,
			OriginalPath: filePath,
			FileSize:     fileSize,
			MimeType:     header.Header.Get("Content-Type"),
			Status:       "uploaded",
		}

		if err := s.store.File().Create(f); err != nil { //s.store.File undefined (type store.Store has no field or method File)
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusCreated, f)
	}
}

func (s *server) handleFilesList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(ctxKeyUser).(*model.User)

		files, err := s.store.File().FindByUserID(user.ID) //s.store.File undefined (type store.Store has no field or method File)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, files)
	}
}

func (s *server) handleFilesGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		user := r.Context().Value(ctxKeyUser).(*model.User)
		file, err := s.store.File().Find(id)
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		if file.UserID != user.ID {
			s.error(w, r, http.StatusForbidden, errors.New("access denied"))
			return
		}

		s.respond(w, r, http.StatusOK, file)
	}
}

func (s *server) handleFilesDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		user := r.Context().Value(ctxKeyUser).(*model.User)
		file, err := s.store.File().Find(id)
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		if file.UserID != user.ID {
			s.error(w, r, http.StatusForbidden, errors.New("access denied"))
			return
		}

		os.Remove(file.OriginalPath)

		if err := s.store.File().Delete(id); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, map[string]string{"status": "deleted"})
	}
}

// Запрос перевода файла
func (s *server) handleFilesTranslate() http.HandlerFunc {
	type request struct {
		SourceLang string `json:"source_lang"`
		TargetLang string `json:"target_lang"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		fileID, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		user := r.Context().Value(ctxKeyUser).(*model.User)
		file, err := s.store.File().Find(fileID)
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		if file.UserID != user.ID {
			s.error(w, r, http.StatusForbidden, errors.New("access denied"))
			return
		}

		// ПРОВЕРЯЕМ ПОДДЕРЖИВАЕМЫЙ ФОРМАТ
		if !s.fileProcessor.SupportsFile(file.MimeType) {
			s.error(w, r, http.StatusUnprocessableEntity,
				errors.New("file format not supported for translation"))
			return
		}

		// ИЗВЛЕКАЕМ ТЕКСТ ИЗ ФАЙЛА
		text, err := s.fileProcessor.ExtractText(file.OriginalPath)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError,
				fmt.Errorf("failed to read file: %w", err))
			return
		}

		// ВЫПОЛНЯЕМ ПЕРЕВОД
		translatedText, err := s.translator.Translate(text, req.SourceLang, req.TargetLang)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError,
				fmt.Errorf("translation failed: %w", err))
			return
		}

		// СОХРАНЯЕМ ПЕРЕВЕДЕННЫЙ ФАЙЛ
		// СОХРАНЯЕМ ПЕРЕВЕДЕННЫЙ ФАЙЛ
		translatedPath, err := s.fileProcessor.SaveTranslatedText(translatedText, file.OriginalPath)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError,
				fmt.Errorf("failed to save translated file: %w", err))
			return
		}

		// СОЗДАЕМ ЗАПИСЬ О ПЕРЕВОДЕ
		translation := &model.Translation{
			FileID:         file.ID,
			SourceLang:     req.SourceLang,
			TargetLang:     req.TargetLang,
			Status:         "completed",
			TranslatedPath: translatedPath,
			CompletedAt:    time.Now(),
		}

		if err := s.store.Translation().Create(translation); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		// ОБНОВЛЯЕМ СТАТУС ФАЙЛА
		file.Status = "translated"
		if err := s.store.File().Update(file); err != nil {
			s.logger.Warnf("failed to update file status: %v", err)
		}

		s.respond(w, r, http.StatusCreated, translation)
	}
}

func (s *server) handleFilesDownload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		user := r.Context().Value(ctxKeyUser).(*model.User)
		file, err := s.store.File().Find(id)
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		if file.UserID != user.ID {
			s.error(w, r, http.StatusForbidden, errors.New("access denied"))
			return
		}

		w.Header().Set("Content-Disposition", "attachment; filename="+file.Filename)
		w.Header().Set("Content-Type", file.MimeType)
		http.ServeFile(w, r, file.OriginalPath)
	}
}

func (s *server) handleTranslationsList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(ctxKeyUser).(*model.User)

		translations, err := s.store.Translation().FindByUserID(user.ID)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, translations)
	}
}

func (s *server) handleTranslationsGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		user := r.Context().Value(ctxKeyUser).(*model.User)
		translation, err := s.store.Translation().Find(id)
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		file, err := s.store.File().Find(translation.FileID)
		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		if file.UserID != user.ID {
			s.error(w, r, http.StatusForbidden, errors.New("access denied"))
			return
		}

		s.respond(w, r, http.StatusOK, translation)
	}
}
