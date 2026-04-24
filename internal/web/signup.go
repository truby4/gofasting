package web

import (
	"errors"
	"net/http"

	"github.com/truby4/go-fasting/internal/auth"
)

func (h *Handler) signup(w http.ResponseWriter, r *http.Request) {
	data := h.newTemplateData(r)
	data.Form = auth.RegistrationForm{}
	h.render(w, r, http.StatusOK, "signup.tmpl", data)
}

func (h *Handler) signupPost(w http.ResponseWriter, r *http.Request) {
	var f auth.RegistrationForm
	data := h.newTemplateData(r)

	err := h.decodePostForm(r, &f)
	if err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	err = h.auth.Registration(&f)
	if err != nil {
		if errors.Is(err, auth.ErrValidation) || !f.Valid() {
			data.Form = f
			h.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
			return
		} else {
			h.serverError(w, r, err)
		}
		return
	}
	http.Redirect(w, r, "/signin", http.StatusSeeOther)
}
