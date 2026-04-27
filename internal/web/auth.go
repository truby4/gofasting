package web

import (
	"errors"
	"net/http"

	"github.com/truby4/gofasting/internal/auth"
)

type authForm struct {
	Email    string `form:"email"`
	Password string `form:"password"`

	Validator
}

func (h *Handler) signin(w http.ResponseWriter, r *http.Request) {
	data := h.newTemplateData(r)
	data.Form = authForm{}
	h.render(w, r, http.StatusOK, "signin.tmpl", data)
}

func (h *Handler) signinPost(w http.ResponseWriter, r *http.Request) {
	var form authForm

	err := h.decodePostForm(r, &form)
	if err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	id, err := h.auth.Authenticate(form.Email, form.Password)
	if err != nil {
		var verrs auth.ValidationError

		switch {
		case errors.As(err, &verrs):
			form.AddFieldErrors(verrs.Fields)

		case errors.Is(err, auth.ErrInvalidCredentials):
			form.AddNonFieldError("Wrong email or password")

		default:
			h.serverError(w, r, err)
			return
		}

		data := h.newTemplateData(r)
		data.Form = form
		h.render(w, r, http.StatusUnprocessableEntity, "signin.tmpl", data)
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

func (h *Handler) signup(w http.ResponseWriter, r *http.Request) {
	data := h.newTemplateData(r)
	data.Form = authForm{}
	h.render(w, r, http.StatusOK, "signup.tmpl", data)
}

func (h *Handler) signupPost(w http.ResponseWriter, r *http.Request) {
	var form authForm

	err := h.decodePostForm(r, &form)
	if err != nil {
		h.clientError(w, http.StatusBadRequest)
		return
	}

	err = h.auth.Register(form.Email, form.Password)
	if err != nil {
		var verrs auth.ValidationError

		switch {
		case errors.Is(err, auth.ErrDuplicateEmail):
			form.AddFieldError("email", "Email address is already in use")

		case errors.As(err, &verrs):
			form.AddFieldErrors(verrs.Fields)

		default:
			h.serverError(w, r, err)
			return
		}

		data := h.newTemplateData(r)
		data.Form = form
		h.render(w, r, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	http.Redirect(w, r, "/signin", http.StatusSeeOther)
}

func (h *Handler) signoutPost(w http.ResponseWriter, r *http.Request) {
	err := h.sessionManager.RenewToken(r.Context())
	if err != nil {
		h.serverError(w, r, err)
		return
	}

	h.sessionManager.Remove(r.Context(), "authenticatedUserID")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
