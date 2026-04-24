package web

import (
	"html/template"

	"github.com/charmbracelet/log"
	"github.com/truby4/go-fasting/internal/store"
)

type Handler struct {
	logger        *log.Logger
	templateCache map[string]*template.Template
	store         *store.Store
}

func NewHandler(logger *log.Logger, store *store.Store) (*Handler, error) {
	tc, err := newTemplateCache()
	if err != nil {
		return nil, err
	}
	return &Handler{
		logger:        logger.WithPrefix("WEB"),
		templateCache: tc,
		store:         store,
	}, nil
}
