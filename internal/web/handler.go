package web

import (
	"html/template"

	"github.com/charmbracelet/log"
)

type Handler struct {
	logger        *log.Logger
	templateCache map[string]*template.Template
}

func NewHandler(logger *log.Logger) *Handler {
	tc, err := newTemplateCache()
	if err != nil {
		log.Fatal(err.Error())
	}
	return &Handler{logger: logger, templateCache: tc}
}
