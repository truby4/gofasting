package api

import "log/slog"

type Handler struct {
	logger *slog.Logger
}

func NewHandler(logger *slog.Logger) (*Handler, error) {
	return &Handler{
		logger: logger.With("component", "API"),
	}, nil
}
