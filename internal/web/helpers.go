package web

import (
	"errors"
	"net/http"

	"github.com/go-playground/form/v4"
)

func (h *Handler) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = h.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		if _, ok := errors.AsType[*form.InvalidDecoderError](err); ok {
			panic(err)
		}

		return err
	}

	return nil
}

func (h *Handler) isAuthenticated(r *http.Request) bool {
	return h.sessionManager.Exists(r.Context(), "authenticatedUserID")
}
