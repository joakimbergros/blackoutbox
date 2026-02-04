package systems

import (
	"blackoutbox/internal/models"
	"blackoutbox/internal/stores"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type SystemHandler struct {
	SystemStore stores.SystemStoreInterface
	UploadRoot  string // e.g. "uploads"
}

// Sync replaces all documents and files for a system.
// Expected payload:
//
//	{
//	  "system_id": "system-123",
//	  "documents": [
//	    {
//	      "file_id": "file-1",
//	      "file_path": "uploads/system-123/file1.pdf",
//	      "print_at": 1710000000,
//	      "tags": ["invoice", "pdf"]
//	    }
//	  ]
//	}
func (h *SystemHandler) Sync() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		systemId := r.PathValue("system_id")
		if systemId == "" {
			http.Error(w, "system_id is required", http.StatusBadRequest)
			return
		}

		var documents []models.Document
		if err := json.NewDecoder(r.Body).Decode(&documents); err != nil {
			http.Error(w, "invalid JSON payload", http.StatusBadRequest)
			return
		}

		now := time.Now()

		for i := range documents {
			documents[i].SystemId = systemId
			documents[i].UpdatedAt = &now
			documents[i].DeletedAt = nil
		}

		if err := h.SystemStore.Sync(systemId, documents); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// Delete removes all documents and files for a system.
func (h *SystemHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		systemId := r.PathValue("system_id")
		if systemId == "" {
			http.Error(w, "system_id is required", http.StatusBadRequest)
			return
		}

		if err := h.SystemStore.DeleteSystem(systemId); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
