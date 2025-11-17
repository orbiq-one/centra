package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cheetahbyte/centra/internal/config"
	"github.com/go-chi/chi/v5"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func handleContent(w http.ResponseWriter, r *http.Request) {
	path := chi.URLParam(r, "*")
	path = strings.Trim(path, "/")

	if path == "" {
		writeJSON(w, http.StatusBadRequest, "Invalid content path")
		return
	}
	parts := strings.SplitN(path, "/", 2)
	collection := parts[0]

	if len(parts) == 1 {
		items, err := config.GetCollection(collection)
		if err != nil {
			writeJSON(w, 500, map[string]any{
				"error": err.Error(),
			})
			return
		}

		writeJSON(w, 200, map[string]any{
			"collection": collection,
			"items":      items,
		})
	} else {
		slug := parts[1]
		entry, err := config.GetEntry(collection, slug)
		if err != nil {
			if err == config.ErrNotFound {
				writeJSON(w, 404, map[string]any{
					"error":      "Not found",
					"collection": collection,
					"slug":       slug,
				})
				return
			}
			writeJSON(w, 500, map[string]any{
				"error": err.Error(),
			})
			return
		}
		writeJSON(w, 200, entry)
	}
}
