package statuses

import (
	"encoding/json"
	"net/http"

	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/auth"
)

type Media struct {
	MediaId     int64  `json:"media_id"`
	Description string `json:"description"`
	Url string `json:"url"`
}

type AddRequest struct {
	Status string `json:"status"`
	Medias []Media `json:"medias"`
}

func (h *handler) Post(w http.ResponseWriter, r *http.Request) {
	//ctx := r.Context()

	var req AddRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	account := auth.AccountOf(r)
	media := object.MediaAttachment{ID: req.Medias[0].MediaId, Description: req.Medias[0].Description, Url: req.Medias[0].Url}
	status := object.Status{}
	status.Account = *account
	status.Content = req.Status
	status.MediaAttachment = media

	
	LastInsertId, err := h.sr.PostStatus(status); 
	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

	status.ID = LastInsertId

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
