package printjobs

import (
	"blackoutbox/internal/response"
	"blackoutbox/internal/stores"
	"encoding/json"
	"net/http"
	"strconv"
)

type PrintJobHandler struct {
	Store stores.PrintJobStoreInterface
}

func (h *PrintJobHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jobs, err := h.Store.Get()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data, err := json.Marshal(jobs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response.JSON(w, http.StatusOK, data)
	}
}

func (h *PrintJobHandler) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "id is required", http.StatusBadRequest)
			return
		}

		job, err := h.Store.GetById(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		data, err := json.Marshal(job)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response.JSON(w, http.StatusOK, data)
	}
}

func (h *PrintJobHandler) GetStuck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		thresholdStr := r.URL.Query().Get("threshold")
		thresholdSeconds := 300

		if thresholdStr != "" {
			threshold, err := strconv.Atoi(thresholdStr)
			if err != nil {
				http.Error(w, "threshold must be a valid integer", http.StatusBadRequest)
				return
			}
			thresholdSeconds = threshold
		}

		jobs, err := h.Store.GetStuckJobs(thresholdSeconds)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data, err := json.Marshal(jobs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response.JSON(w, http.StatusOK, data)
	}
}
