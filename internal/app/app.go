package app

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/truby4/gofasting/internal/api"
	"github.com/truby4/gofasting/internal/auth"
	"github.com/truby4/gofasting/internal/web"
)

type Application struct {
	Logger         *slog.Logger
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
	db             *sql.DB
	auth           *auth.Service

	web *web.Handler
	api *api.Handler
}

func New() (*Application, error) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB()
	if err != nil {
		return nil, err
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	auth := auth.New(db)

	web, err := web.NewHandler(logger, formDecoder, auth, sessionManager)
	if err != nil {
		return nil, err
	}

	api, err := api.NewHandler(logger)
	if err != nil {
		return nil, err
	}

	return &Application{
		Logger:         logger,
		web:            web,
		api:            api,
		sessionManager: sessionManager,
		db:             db,
	}, nil
}

func (app *Application) Serve(addr string) error {
	app.Logger.Info("Starting server", "addr", addr)
	srv := &http.Server{
		Addr:         addr,
		Handler:      app.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err := srv.ListenAndServe()
	return err
}

func (app *Application) Close() error {
	return app.db.Close()
}
