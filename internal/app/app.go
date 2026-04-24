package app

import (
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/charmbracelet/log"
	"github.com/go-playground/form/v4"
	"github.com/truby4/go-fasting/internal/api"
	"github.com/truby4/go-fasting/internal/auth"
	"github.com/truby4/go-fasting/internal/store"
	"github.com/truby4/go-fasting/internal/web"
)

type Application struct {
	Logger         *log.Logger
	store          *store.Store
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager

	auth *auth.Service

	web *web.Handler
	api *api.Handler
}

func New() (*Application, error) {
	logger := log.NewWithOptions(os.Stderr, log.Options{ReportTimestamp: true})

	store, err := store.New(logger)
	if err != nil {
		return nil, err
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(store.DB)
	sessionManager.Lifetime = 12 * time.Hour

	auth, err := auth.New(logger, store.Users)
	if err != nil {
		store.Close()
		return nil, err
	}

	web, err := web.NewHandler(logger, formDecoder, &auth, sessionManager)
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
		Logger:         logger,
		web:            web,
		api:            api,
		store:          store,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
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
