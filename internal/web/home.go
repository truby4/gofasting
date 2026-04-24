package web

import (
	"net/http"
)

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	data := h.newTemplateData(r)
	h.render(w, r, http.StatusOK, "home.tmpl", data)
}
