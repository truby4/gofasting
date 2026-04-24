package web

import (
	"html/template"

	"github.com/alexedwards/scs/v2"
	"github.com/charmbracelet/log"
	"github.com/go-playground/form/v4"
	"github.com/truby4/go-fasting/internal/auth"
)

type Handler struct {
	logger         *log.Logger
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	auth           *auth.Service
	sessionManager *scs.SessionManager
}

func NewHandler(
	logger *log.Logger,
	formDecoder *form.Decoder,
	auth *auth.Service,
	sessionManager *scs.SessionManager,
) (*Handler, error) {
	tc, err := newTemplateCache()
	if err != nil {
		return nil, err
	}
	return &Handler{
		logger:         logger.WithPrefix("WEB"),
		templateCache:  tc,
		formDecoder:    formDecoder,
		auth:           auth,
		sessionManager: sessionManager,
	}, nil
}
