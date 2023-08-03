package statuses

import (
	"encoding/json"
	"net/http"
	"github.com/go-chi/chi/v5"
	"strconv"
)


func (h *handler) FindStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idstr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid id in URL path", http.StatusBadRequest)
		return
	}

	status, AccountId, err := h.sr.FindById(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	account, err := h.ar.FindById(ctx, AccountId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	status.Account = *account

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
