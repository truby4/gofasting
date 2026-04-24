package web

import (
	"net/http"
)

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, http.StatusOK, "home.tmpl", templateData{})
}
