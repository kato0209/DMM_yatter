package statuses

import (
	"net/http"
	"yatter-backend-go/app/domain/repository"
	"yatter-backend-go/app/handler/auth"
	"fmt"

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
	fmt.Println(99)
	r.With(auth.Middleware(ar)).Post("/", h.Post)
	fmt.Println(100)
	return r
}
