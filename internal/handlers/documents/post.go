package documents

import (
	"blackoutbox/internal/models"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func (h *DocumentHandler) Post() http.HandlerFunc {
	const maxFileSize = 10 << 20 // 10MB

	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(maxFileSize); err != nil {
			http.Error(w, "Unable to parse multipart form", http.StatusBadRequest)
			return
		}

		externalId := r.FormValue("ext_id")
		if externalId == "" {
			http.Error(w, "ext_id is required", http.StatusBadRequest)
			return
		}

		isSystem := r.FormValue("is_system")
		if isSystem == "" {
			http.Error(w, "is_system is required", http.StatusBadRequest)
			return
		}

		systemBool, err := strconv.ParseBool(isSystem)
		if err != nil {
			http.Error(w, "is_system needs to be a boolean value", http.StatusBadRequest)
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

		uploadDir := filepath.Join("uploads", externalId)
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

		now := time.Now()

		if err := h.Store.Add(models.Document{
			ExternalId: externalId,
			IsSystem:   systemBool,
			FilePath:   filepath.Join("uploads", externalId, filename),
			UpdatedAt:  &now,
			DeletedAt:  nil,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
