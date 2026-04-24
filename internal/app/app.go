package app

import (
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/truby4/go-fasting/internal/api"
	"github.com/truby4/go-fasting/internal/store"
	"github.com/truby4/go-fasting/internal/web"
)

type Application struct {
	Logger *log.Logger
	store  *store.Store

	web *web.Handler
	api *api.Handler
}

func New() (*Application, error) {
	logger := log.NewWithOptions(os.Stderr, log.Options{ReportTimestamp: true})

	store, err := store.New(logger)
	if err != nil {
		return nil, err
	}

	web, err := web.NewHandler(logger, store)
	if err != nil {
		store.Close()
		return nil, err
	}

	api, err := api.NewHandler(logger, store)
	if err != nil {
		store.Close()
		return nil, err
	}

	return &Application{
		Logger: logger,
		web:    web,
		api:    api,
		store:  store,
	}, nil
}

func (app *Application) Serve(addr string) error {
	app.Logger.Infof("Starting server on localhost %s", addr)
	srv := &http.Server{
		Addr:         addr,
		Handler:      app.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     app.Logger.StandardLog(log.StandardLogOptions{ForceLevel: log.ErrorLevel}),
	}

	err := srv.ListenAndServe()
	return err
}

func (app *Application) Close() error {
	return app.store.Close()
}
