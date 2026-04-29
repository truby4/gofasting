package web

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(preventCSRF)

	r.With(h.authMiddleware).Get("/", h.home)

	r.Get("/signin", h.signin)
	r.Post("/signin", h.signinPost)

	r.Get("/signup", h.signup)
	r.Post("/signup", h.signupPost)

	r.Post("/signout", h.signoutPost)

	r.With(h.authMiddleware).Post("/fast/start", h.fastStartPost)
	r.With(h.authMiddleware).Post("/fast/end", h.fastEndPost)

	return r
}
