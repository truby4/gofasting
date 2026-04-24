package web

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/", h.home)

	return r
}
