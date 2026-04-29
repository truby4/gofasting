package web

import (
	"html/template"
	"log/slog"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/truby4/gofasting/internal/auth"
	"github.com/truby4/gofasting/internal/fasts"
)

type Handler struct {
	logger         *slog.Logger
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	auth           *auth.Service
	fasts          *fasts.Store
	sessionManager *scs.SessionManager
}

func NewHandler(
	logger *slog.Logger,
	auth *auth.Service,
	fasts *fasts.Store,
	sessionManager *scs.SessionManager,
) (*Handler, error) {
	tc, err := newTemplateCache()
	if err != nil {
		return nil, err
	}
	formDecoder := form.NewDecoder()
	return &Handler{
		logger:         logger.With("component", "WEB"),
		templateCache:  tc,
		formDecoder:    formDecoder,
		auth:           auth,
		fasts:          fasts,
		sessionManager: sessionManager,
	}, nil
}
