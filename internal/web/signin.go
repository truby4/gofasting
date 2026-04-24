package web

import (
	"errors"
	"net/http"

	"github.com/truby4/go-fasting/internal/auth"
)

func (h *Handler) signin(w http.ResponseWriter, r *http.Request) {
	data := h.newTemplateData(r)
	data.Form = auth.AuthenticationForm{}
	h.render(w, r, http.StatusOK, "signin.tmpl", data)
}

func (h *Handler) signinPost(w http.ResponseWriter, r *http.Request) {
	var f auth.AuthenticationForm
	data := h.newTemplateData(r)

	err := h.decodePostForm(r, &f)
	if err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	id, err := h.auth.Authenticate(&f)
	if err != nil {
		if errors.Is(err, auth.ErrValidation) || !f.Valid() || errors.Is(err, auth.ErrInvalidCredentials) {
			data.Form = f
			h.render(w, r, http.StatusUnprocessableEntity, "signin.tmpl", data)
			return
		} else {
			h.serverError(w, r, err)
		}
		return
	}

	err = h.sessionManager.RenewToken(r.Context())
	if err != nil {
		h.serverError(w, r, err)
		return
	}

	h.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
