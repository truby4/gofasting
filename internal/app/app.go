package app

import (
	"crypto/tls"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/truby4/gofasting/internal/api"
	"github.com/truby4/gofasting/internal/auth"
	"github.com/truby4/gofasting/internal/fasts"
	"github.com/truby4/gofasting/internal/web"
)

type Application struct {
	Logger         *slog.Logger
	sessionManager *scs.SessionManager
	db             *sql.DB

	web *web.Handler
	api *api.Handler
}

func New(level slog.Level) (*Application, error) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))

	db, err := openDB()
	if err != nil {
		return nil, err
	}

	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	auth := auth.New(db)
	fasts := fasts.New(db)

	web, err := web.NewHandler(
		logger,
		auth,
		fasts,
		sessionManager,
	)
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

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:         addr,
		Handler:      app.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(app.Logger.Handler(), slog.LevelError),
		TLSConfig:    tlsConfig,
	}

	return srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
}

func (app *Application) Close() error {
	return app.db.Close()
}
