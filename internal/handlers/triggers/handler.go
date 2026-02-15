// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package triggers

import (
	"blackoutbox/internal/models"
	"blackoutbox/internal/response"
	"blackoutbox/internal/stores"
	"blackoutbox/internal/validation"
	"encoding/json"
	"net/http"
	"time"
)

type TriggerHandler struct {
	Store stores.TriggerStoreInterface
}

func (h *TriggerHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		triggers, err := h.Store.Get()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data, err := json.Marshal(triggers)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response.JSON(w, http.StatusOK, data)
	}
}

func (h *TriggerHandler) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "id is required", http.StatusBadRequest)
			return
		}

		trigger, err := h.Store.GetById(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		data, err := json.Marshal(trigger)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response.JSON(w, http.StatusOK, data)
	}
}

func (h *TriggerHandler) Post() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			SystemId      string `json:"system_id"`
			Url           string `json:"url"`
			BufferSeconds *int   `json:"buffer_seconds"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.SystemId == "" {
			http.Error(w, "system_id is required", http.StatusBadRequest)
			return
		}

		if req.Url == "" {
			http.Error(w, "url is required", http.StatusBadRequest)
			return
		}

		//TODO Add more validation for SSRF attacks
		if err := validation.ValidateTriggerURL(req.Url); err != nil {
			http.Error(w, "url is not public", http.StatusBadRequest)
			return
		}

		bufferSeconds := 300
		if req.BufferSeconds != nil {
			bufferSeconds = *req.BufferSeconds
		}

		now := time.Now().Unix()

		trigger := models.Trigger{
			SystemId:      req.SystemId,
			Url:           req.Url,
			BufferSeconds: bufferSeconds,
			Status:        "ok",
			RetryCount:    0,
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		if err := h.Store.Add(trigger); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func (h *TriggerHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "id is required", http.StatusBadRequest)
			return
		}

		if _, err := h.Store.GetById(id); err != nil {
			http.Error(w, "Trigger not found", http.StatusNotFound)
			return
		}

		if err := h.Store.Delete(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
