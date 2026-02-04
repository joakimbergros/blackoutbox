package documents

import (
	"blackoutbox/internal/models"
	"blackoutbox/internal/response"
	"blackoutbox/internal/stores"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type DocumentHandler struct {
	Store stores.DocumentStoreInterface
}

func (h *DocumentHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		systemIdFilter := r.URL.Query().Get("system-id")
		fileIdFilter := r.URL.Query().Get("file-id")

		var data []byte
		var marshallErr error

		if systemIdFilter != "" {
			documents, err := h.Store.GetBySystemId(systemIdFilter)
			if err != nil {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			data, marshallErr = json.Marshal(documents)
		} else if fileIdFilter != "" {
			document, err := h.Store.GetByFileId(fileIdFilter)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			data, marshallErr = json.Marshal(document)
		} else {
			documents, err := h.Store.Get()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			data, marshallErr = json.Marshal(documents)
		}

		if marshallErr != nil {
			http.Error(w, marshallErr.Error(), http.StatusInternalServerError)
			return
		}

		response.JSON(w, http.StatusOK, data)
	}
}

func (h *DocumentHandler) Post() http.HandlerFunc {
	const maxFileSize = 10 << 20 // 10MB

	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(maxFileSize); err != nil {
			http.Error(w, "Unable to parse multipart form", http.StatusBadRequest)
			return
		}

		systemId := r.FormValue("system_id")
		if systemId == "" {
			http.Error(w, "system_id is required", http.StatusBadRequest)
			return
		}

		fileId := r.FormValue("file_id")
		if fileId == "" {
			http.Error(w, "file_id is required", http.StatusBadRequest)
			return
		}

		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "file is required", http.StatusBadRequest)
			return
		}
		defer file.Close()

		if fileHeader.Size > maxFileSize {
			http.Error(w, "File size exceeds 10MB limit", http.StatusBadRequest)
			return
		}

		uploadDir := filepath.Join("uploads", systemId)
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			http.Error(w, "Failed to create upload directory", http.StatusInternalServerError)
			return
		}

		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), fileHeader.Filename)
		filePath := filepath.Join(uploadDir, filename)

		dst, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Failed to create file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			return
		}

		var printAt *int64
		printAtStr := r.FormValue("print_at")
		if printAtStr != "" {
			timestamp, err := strconv.ParseInt(printAtStr, 10, 64)
			if err != nil {
				http.Error(w, "print_at must be a valid Unix timestamp", http.StatusBadRequest)
				return
			}
			printAt = &timestamp
		}

		var tags []string
		tagsStr := r.FormValue("tags")
		if tagsStr != "" {
			if err := json.Unmarshal([]byte(tagsStr), &tags); err != nil {
				http.Error(w, "tags must be a valid JSON array", http.StatusBadRequest)
				return
			}
		}

		now := time.Now()

		if err := h.Store.Add(models.Document{
			SystemId:      systemId,
			FileId:        fileId,
			FilePath:      filepath.Join("uploads", systemId, filename),
			PrintAt:       printAt,
			LastPrintedAt: nil,
			Tags:          tags,
			UpdatedAt:     &now,
			DeletedAt:     nil,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func (h *DocumentHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (h *DocumentHandler) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "id is required", http.StatusBadRequest)
			return
		}

		document, err := h.Store.GetById(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response.JSON(w, http.StatusOK, document)
	}
}
