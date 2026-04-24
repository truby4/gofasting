package auth

import (
	"github.com/charmbracelet/log"
	"github.com/truby4/go-fasting/internal/store"
)

type Service struct {
	logger *log.Logger
	store  *store.UserStore
}

func New(logger *log.Logger, store *store.UserStore) (Service, error) {
	return Service{
		logger: logger,
		store:  store,
	}, nil
}
