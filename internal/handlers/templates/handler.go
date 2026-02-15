// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package templates

import (
	"blackoutbox/internal/models"
	"blackoutbox/internal/response"
	"blackoutbox/internal/stores"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type TemplatesHandler struct {
	Store stores.TemplateStoreInterface
}

func (h *TemplatesHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		systemIdFilter := r.URL.Query().Get("system-id")
		fileReference := r.URL.Query().Get("file-id")

		systemIntId, err := strconv.ParseInt(systemIdFilter, 10, 64)
		if err != nil {
			return
		}

		var data any

		if systemIdFilter != "" {
			documents, err := h.Store.GetBySystemId(systemIntId)
			if err != nil {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			data = documents
		} else if fileReference != "" {
			document, err := h.Store.GetByFileReference(fileReference)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			data = document
		} else {
			documents, err := h.Store.Get()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			data = documents
		}

		response.JSON(w, http.StatusOK, data)
	}
}

func (h *TemplatesHandler) Post() http.HandlerFunc {
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

		systemIntId, err := strconv.ParseInt(systemId, 10, 64)
		if err != nil {
			return
		}

		if strings.Contains(systemId, "../") || strings.Contains(systemId, "..\\") {
			http.Error(w, "Invalid file path", http.StatusBadRequest)
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

		uploadDir := filepath.Join("templates", fileId)
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

		now := time.Now().Unix()

		if err := h.Store.Add(models.Template{
			SystemId:      systemIntId,
			FileReference: fileId,
			FilePath:      filePath,
			Description:   r.FormValue("description"),
			CreatedAt:     now,
			DeletedAt:     nil,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func (h *TemplatesHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
