package timelines

import (
	"net/http"
	"yatter-backend-go/app/domain/repository"

	"github.com/go-chi/chi/v5"
)

// Implementation of handler
type handler struct {
	sr repository.Status
	ar repository.Account
}

// Create Handler for `/v1/accounts/`
func NewRouter(sr repository.Status, ar repository.Account) http.Handler {
	r := chi.NewRouter()

	h := &handler{sr,ar}
	r.Get("/public", h.PublicTimelines)
	return r
}
