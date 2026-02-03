package documents

import (
	"encoding/json"
	"net/http"
)

func (h *DocumentHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		documents, err := h.Store.Get()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json, err := json.Marshal(documents)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(json)
		return
	}
}
