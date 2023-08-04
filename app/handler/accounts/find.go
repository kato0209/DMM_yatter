package accounts

import (
	"encoding/json"
	"net/http"
	"github.com/go-chi/chi/v5"
)


// Handle request for `GET /v1/accounts/{username}`
func (h *handler) GetAccountHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	username := chi.URLParam(r, "username")
	if username == "" {
		http.Error(w, "username not found in URL path", http.StatusBadRequest)
		return
	}

	account, err := h.ar.FindByUsername(ctx, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(account); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
