// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package systems

import (
	"blackoutbox/internal/models"
	"blackoutbox/internal/stores"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type SystemHandler struct {
	SystemStore stores.SystemStoreInterface
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
		systemRef := r.PathValue("system_id")
		if systemRef == "" {
			http.Error(w, "system_id is required", http.StatusBadRequest)
			return
		}

		// Get the system by reference to get the internal ID
		system, err := h.SystemStore.GetSystemByReference(systemRef)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if system == nil {
			http.Error(w, "System not found", http.StatusNotFound)
			return
		}

		var documents []models.Document
		if err := json.NewDecoder(r.Body).Decode(&documents); err != nil {
			http.Error(w, "invalid JSON payload", http.StatusBadRequest)
			return
		}

		now := time.Now().Unix()

		for i := range documents {
			documents[i].SystemId = system.Id
			documents[i].UpdatedAt = &now
			documents[i].DeletedAt = nil
		}

		if err := h.SystemStore.Sync(system.Id, documents); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// Delete removes all documents and files for a system.
func (h *SystemHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		systemRef := r.PathValue("id")
		if systemRef == "" {
			http.Error(w, "system_id is required", http.StatusBadRequest)
			return
		}

		// Get the system by reference to get the internal ID
		system, err := h.SystemStore.GetSystemByReference(systemRef)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if system == nil {
			http.Error(w, "System not found", http.StatusNotFound)
			return
		}

		if err := h.SystemStore.DeleteSystem(system.Id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// GetSystems handles GET /systems - List all systems
func (h *SystemHandler) GetSystems() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		systems, err := h.SystemStore.GetSystems()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(systems)
	}
}

// GetSystem handles GET /systems/{id} - Get a specific system
func (h *SystemHandler) GetSystem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reference := r.PathValue("id")
		if reference == "" {
			http.Error(w, "System reference is required", http.StatusBadRequest)
			return
		}

		// Try to get system by reference first (user-friendly)
		system, err := h.SystemStore.GetSystemByReference(reference)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if system == nil {
			// Fall back to getting by ID (for internal use)
			// This requires parsing the reference as an integer
			id, parseErr := strconv.ParseInt(reference, 10, 64)
			if parseErr != nil {
				http.Error(w, "System not found", http.StatusNotFound)
				return
			}
			system, err = h.SystemStore.GetSystemById(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		if system == nil {
			http.Error(w, "System not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(system)
	}
}

// CreateSystem handles POST /systems - Create a new system
func (h *SystemHandler) CreateSystem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var system models.System
		err := json.NewDecoder(r.Body).Decode(&system)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Validate required fields
		if system.Reference == "" {
			http.Error(w, "System reference is required", http.StatusBadRequest)
			return
		}

		if system.Name == "" {
			http.Error(w, "System name is required", http.StatusBadRequest)
			return
		}

		// Set timestamps
		now := time.Now().Unix()
		system.CreatedAt = now
		system.UpdatedAt = now

		err = h.SystemStore.AddSystem(system)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Get the created system to return with the generated ID
		createdSystem, err := h.SystemStore.GetSystemByReference(system.Reference)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdSystem)
	}
}

// UpdateSystem handles PUT /systems/{id} - Update a system
func (h *SystemHandler) UpdateSystem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reference := r.PathValue("id")
		if reference == "" {
			http.Error(w, "System reference is required", http.StatusBadRequest)
			return
		}

		// Get the existing system to get the internal ID
		existingSystem, err := h.SystemStore.GetSystemByReference(reference)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if existingSystem == nil {
			http.Error(w, "System not found", http.StatusNotFound)
			return
		}

		var updatedSystem models.System
		err = json.NewDecoder(r.Body).Decode(&updatedSystem)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Ensure the reference matches the path (if provided in body)
		if updatedSystem.Reference != "" && updatedSystem.Reference != reference {
			http.Error(w, "System reference mismatch", http.StatusBadRequest)
			return
		}

		// Validate required fields
		if updatedSystem.Name == "" {
			http.Error(w, "System name is required", http.StatusBadRequest)
			return
		}

		// Set the ID and reference for update
		updatedSystem.Id = existingSystem.Id
		if updatedSystem.Reference == "" {
			updatedSystem.Reference = existingSystem.Reference
		}

		err = h.SystemStore.UpdateSystem(updatedSystem)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedSystem)
	}
}

// DeleteSystem handles DELETE /systems/{id} - Delete a system (soft delete)
func (h *SystemHandler) DeleteSystem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reference := r.PathValue("id")
		if reference == "" {
			http.Error(w, "System reference is required", http.StatusBadRequest)
			return
		}

		// Get the system by reference to get the internal ID
		system, err := h.SystemStore.GetSystemByReference(reference)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if system == nil {
			http.Error(w, "System not found", http.StatusNotFound)
			return
		}

		err = h.SystemStore.DeleteSystem(system.Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
