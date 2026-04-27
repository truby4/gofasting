package web

import (
	"html/template"
	"log/slog"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/truby4/gofasting/internal/auth"
)

type Handler struct {
	logger         *slog.Logger
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	auth           *auth.Service
	sessionManager *scs.SessionManager
}

func NewHandler(
	logger *slog.Logger,
	formDecoder *form.Decoder,
	auth *auth.Service,
	sessionManager *scs.SessionManager,
) (*Handler, error) {
	tc, err := newTemplateCache()
	if err != nil {
		return nil, err
	}
	return &Handler{
		logger:         logger.With("component", "WEB"),
		templateCache:  tc,
		formDecoder:    formDecoder,
		auth:           auth,
		sessionManager: sessionManager,
	}, nil
}
