package api

import (
	"github.com/charmbracelet/log"
	"github.com/truby4/go-fasting/internal/store"
)

type Handler struct {
	logger *log.Logger
	store  *store.Store
}

func NewHandler(logger *log.Logger, store *store.Store) (*Handler, error) {
	return &Handler{
		logger: logger.WithPrefix("API"),
		store:  store,
	}, nil
}
